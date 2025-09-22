document.addEventListener('DOMContentLoaded', () => {
    const topicListElement = document.getElementById('topic-list');

    // 模拟从后端获取的话题列表
    const topics = [
        "今天你摸鱼了吗？",
        "聊聊最近看的电影",
        "下班/放学去哪玩",
        "周末爬山搭子",
        "随便聊聊"
    ];

    function displayTopics() {
        // 清空加载提示
        topicListElement.innerHTML = '';

        topics.forEach(topic => {
            const topicItem = document.createElement('div');
            topicItem.className = 'topic-item';
            topicItem.textContent = topic;
            topicItem.addEventListener('click', () => {
                // 1. 生成一个唯一的房间ID
                const roomId = generateUUID();
                
                // 2. 跳转到聊天页面，并带上房间ID和话题作为参数
                // 使用 encodeURIComponent 来确保特殊字符能被正确处理
                window.location.href = `chat.html?room=${roomId}&topic=${encodeURIComponent(topic)}`;
            });
            topicListElement.appendChild(topicItem);
        });
    }

    /**
     * 生成一个简单的唯一标识符 (UUID v4)
     * 在实际应用中，房间ID最好由后端生成以保证绝对唯一
     */
    function generateUUID() {
        return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
            (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
        );
    }

    // 运行
    displayTopics();
});
