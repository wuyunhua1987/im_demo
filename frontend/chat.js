import IMSDK from '@enigma-im/jssdk';
import FingerprintJS from '@fingerprintjs/fingerprintjs';

async function getVisitorId() {
  const fp = await FingerprintJS.load();
  const result = await fp.get();
  return result.visitorId;
}

const baseUrl = "http://localhost:8089"

document.addEventListener('DOMContentLoaded', async () => {
  // 从 DOM 获取元素
  const roomTitleElement = document.getElementById('room-title');
  const messageContainer = document.getElementById('message-container');
  const messageForm = document.getElementById('message-form');
  const messageInput = document.getElementById('message-input');
  const fp = await getVisitorId();
  console.log('获取到的用户标识:', fp);
  const userinfo = await (await fetch(`${baseUrl}/v1/get_info`, {
    method: 'GET',
    headers: {
      'X-Finger-Print': fp
    }
  })).json();
  console.log('获取到的用户令牌:', userinfo.token);

  // 1. 从 URL 解析参数
  const params = new URLSearchParams(window.location.search);
  const roomId = params.get('room');
  const topic = params.get('topic');

  // 如果没有房间ID或话题，则跳转回首页
  if (!roomId || !topic) {
    window.location.href = 'index.html';
    return;
  }

  // 2. 设置房间标题
  roomTitleElement.textContent = topic;

  // =================================================================
  // TODO: 在这里集成你的 WebSocket JSSDK
  // =================================================================
  function initializeWebSocketSDK() {
    console.log(`准备为用户 ${userinfo.user_id} 初始化 WebSocket 连接...`);
    const client = new IMSDK.IMClient({
      url: 'wss://im.miaowankeji.com/ws', // Using a public echo server for demonstration
      token: userinfo.token,
      autoReconnect: true
    });

    console.log('Client initialized. Connecting...');

    // 2. Add event listeners
    client.on('connect', () => {
      console.log('Successfully connected to server!', 'connect');
      // 加入频道
      client.joinChannel(roomId);
    });

    client.on('priv-channel-message', (message) => {
      const { fromUserId, payload } = message
      if (payload && payload.length > 0) {
        let bytes = new Uint8Array(payload);
        let str = new TextDecoder().decode(bytes);
        console.log(`Received message: ${JSON.stringify(message)}`);
        if (fp === fromUserId) {
          displayMessage('我', str, 'sent');
        } else {
          displayMessage(fromUserId, str, 'received');
        }
      }
    });

    client.on('joined-channel', (message) => {
      console.log(`Joined channel ${JSON.stringify(message)} successfully`, 'info');
      if (fp === message.userId) {
        displayMessage('我', '加入了房间', 'sent');
      } else {
        displayMessage(message.userId, '加入了房间', 'received');
      }
    });

    client.on('left-channel', (message) => {
      console.log(`Left channel ${JSON.stringify(message)} successfully`, 'info');
      displayMessage(message.userId, '离开了房间', 'received');
    });

    client.on('error', (error) => {
      console.log('An error occurred.', 'disconnect');
      console.error(error);
    });

    // 3. Connect
    client.connect();
  }
  // =================================================================

  // 3. 处理消息发送
  messageForm.addEventListener('submit', (e) => {
    e.preventDefault(); // 阻止表单默认的提交行为
    const messageText = messageInput.value.trim();

    if (messageText) {
      // 调用发送消息的函数
      sendMessage(messageText);
      // 清空输入框
      messageInput.value = '';
    }
  });

  /**
   * 发送消息
   * @param {string} text 消息内容
   */
  async function sendMessage(text) {
    // b. 调用第三方的HTTP接口来发送消息
    try {
      const apiUrl = `${baseUrl}/v1/send_channel_msg`;

      const response = await fetch(apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          // 可能需要认证信息，如 'Authorization': 'Bearer YOUR_TOKEN'
          'X-Finger-Print': fp
        },
        body: JSON.stringify({
          room_id: roomId,
          user_id: userinfo.user_id,
          content: text
        })
      });

      if (!response.ok) {
        console.error('消息发送失败:', response.statusText);
        // 可选：给用户一个发送失败的提示
      }
    } catch (error) {
      console.error('网络请求错误:', error);
      // 可选：给用户一个发送失败的提示
    }
  }

  /**
   * 在聊天窗口显示一条消息
   * @param {string} sender 发送者昵称
   * @param {string} text 消息内容
   * @param {'sent' | 'received'} type 消息类型 ('sent' 或 'received')
   */
  function displayMessage(sender, text, type) {
    const bubble = document.createElement('div');
    bubble.className = `message-bubble ${type}`;

    const senderElement = document.createElement('div');
    senderElement.className = 'message-sender';
    senderElement.textContent = sender+":";
    bubble.appendChild(senderElement);

    const textElement = document.createElement('div');
    textElement.textContent = text;
    bubble.appendChild(textElement);

    messageContainer.appendChild(bubble);

    // 滚动到底部，以显示最新消息
    messageContainer.scrollTop = messageContainer.scrollHeight;
  }

  // 页面加载后，初始化WebSocket连接
  initializeWebSocketSDK();
});
