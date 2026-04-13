package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 注册根路径处理器
	http.HandleFunc("/", handleHome)

	// 启动服务器
	port := "8080"
	log.Printf("Server starting on http://localhost:%s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	// 只处理根路径
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// 内嵌的 HTML 页面
	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>今天吃什么？</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
            font-family: 'Segoe UI', 'Microsoft YaHei', sans-serif;
        }

        body {
            background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }

        .container {
            text-align: center;
            background-color: white;
            border-radius: 24px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
            padding: 60px 40px;
            max-width: 800px;
            width: 100%;
            position: relative;
            overflow: hidden;
        }

        .container::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 8px;
            background: linear-gradient(90deg, #ff9a9e 0%, #fad0c4 100%);
        }

        h1 {
            color: #333;
            font-size: 3.2rem;
            margin-bottom: 20px;
            font-weight: 700;
            letter-spacing: -0.5px;
        }

        .subtitle {
            color: #666;
            font-size: 1.4rem;
            margin-bottom: 60px;
            line-height: 1.6;
        }

        .button-container {
            margin-bottom: 70px;
        }

        #eatButton {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 28px 70px;
            font-size: 2.2rem;
            font-weight: 700;
            border-radius: 70px;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
            letter-spacing: 1px;
        }

        #eatButton:hover {
            transform: translateY(-5px);
            box-shadow: 0 15px 40px rgba(102, 126, 234, 0.6);
        }

        #eatButton:active {
            transform: translateY(0);
        }

        .result-container {
            min-height: 180px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            margin-top: 30px;
        }

        #result {
            font-size: 4.5rem;
            font-weight: 900;
            color: #333;
            margin-bottom: 20px;
            opacity: 0;
            transform: scale(0.8);
            transition: opacity 0.5s ease, transform 0.5s ease;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.1);
        }

        #result.show {
            opacity: 1;
            transform: scale(1);
        }

        .food-icon {
            font-size: 3rem;
            margin-bottom: 20px;
            opacity: 0.7;
        }

        .history {
            margin-top: 50px;
            padding-top: 30px;
            border-top: 1px solid #eee;
        }

        .history h3 {
            color: #666;
            font-size: 1.2rem;
            margin-bottom: 15px;
            font-weight: 600;
        }

        #historyList {
            display: flex;
            flex-wrap: wrap;
            justify-content: center;
            gap: 12px;
            max-height: 120px;
            overflow-y: auto;
            padding: 10px;
        }

        .history-item {
            background: #f8f9fa;
            padding: 8px 18px;
            border-radius: 20px;
            font-size: 1.1rem;
            color: #555;
            border: 1px solid #e9ecef;
        }

        footer {
            margin-top: 50px;
            color: #999;
            font-size: 0.9rem;
        }

        @media (max-width: 768px) {
            .container {
                padding: 40px 20px;
            }
            h1 {
                font-size: 2.5rem;
            }
            #eatButton {
                padding: 22px 50px;
                font-size: 1.8rem;
            }
            #result {
                font-size: 3.5rem;
            }
        }

        @media (max-width: 480px) {
            h1 {
                font-size: 2rem;
            }
            .subtitle {
                font-size: 1.1rem;
            }
            #eatButton {
                padding: 18px 40px;
                font-size: 1.5rem;
            }
            #result {
                font-size: 2.8rem;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>今天吃什么？</h1>
        <p class="subtitle">点击下方按钮，让命运决定今天的美食吧！</p>
        
        <div class="button-container">
            <button id="eatButton">今天吃什么</button>
        </div>

        <div class="result-container">
            <div class="food-icon">🍽️</div>
            <div id="result"></div>
        </div>

        <div class="history">
            <h3>历史记录</h3>
            <div id="historyList"></div>
        </div>

        <footer>
            <p>每天都要好好吃饭哦！</p>
        </footer>
    </div>

    <script>
        // 食物选项数组
        const foodOptions = [
            "火锅", "烧烤", "麻辣烫", "西餐", "日料", 
            "家常菜", "拉面", "汉堡", "炸鸡", "饺子",
            "披萨", "寿司", "咖喱饭", "炸酱面", "酸菜鱼",
            "烤鸭", "麻辣香锅", "牛肉面", "炒饭", "沙拉"
        ];

        // 获取DOM元素
        const eatButton = document.getElementById('eatButton');
        const resultElement = document.getElementById('result');
        const historyList = document.getElementById('historyList');
        
        // 历史记录数组
        let history = [];

        // 随机选择食物
        function getRandomFood() {
            const randomIndex = Math.floor(Math.random() * foodOptions.length);
            return foodOptions[randomIndex];
        }

        // 更新历史记录显示
        function updateHistory() {
            historyList.innerHTML = '';
            // 只显示最近8条记录
            const recentHistory = history.slice(-8);
            recentHistory.forEach(food => {
                const item = document.createElement('div');
                item.className = 'history-item';
                item.textContent = food;
                historyList.appendChild(item);
            });
        }

        // 按钮点击事件
        eatButton.addEventListener('click', () => {
            // 获取随机食物
            const food = getRandomFood();
            
            // 隐藏结果元素，准备动画
            resultElement.classList.remove('show');
            
            // 短暂延迟后显示新结果
            setTimeout(() => {
                resultElement.textContent = food;
                resultElement.classList.add('show');
                
                // 添加到历史记录
                history.push(food);
                updateHistory();
            }, 300);
        });

        // 页面加载时显示一条随机结果
        window.addEventListener('load', () => {
            const initialFood = getRandomFood();
            resultElement.textContent = initialFood;
            resultElement.classList.add('show');
            history.push(initialFood);
            updateHistory();
        });
    </script>
</body>
</html>`

	// 写入响应
	fmt.Fprint(w, html)
}