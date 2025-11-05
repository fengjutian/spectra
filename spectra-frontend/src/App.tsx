import { useState } from 'react'
// import reactLogo from './assets/react.svg'
// import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [count, setCount] = useState(0)
  const [response, setResponse] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)

  // 记录错误日志的方法
  const recordErrorLog = async () => {
    setLoading(true)
    setResponse('')
    
    try {
      const errorData = {
        timestamp: new Date().toISOString(),
        project_id: 'demo-project',
        session_id: 'session-123',
        trace_id: 'trace-456',
        user_id: 'user-789',
        url: window.location.href,
        referrer: document.referrer,
        type: 'javascript_error',
        name: 'DemoError',
        message: '这是一个演示错误日志',
        extra: JSON.stringify({ 
          browser: navigator.userAgent,
          count: count 
        })
      }

      const response = await fetch('http://localhost:8080/api/error-logs', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(errorData)
      })

      if (response.ok) {
        const result = await response.json()
        setResponse(`✅ 成功记录错误日志！\n\n响应数据:\n${JSON.stringify(result, null, 2)}`)
      } else {
        const errorText = await response.text()
        setResponse(`❌ 请求失败: ${response.status} ${response.statusText}\n\n错误详情:\n${errorText}`)
      }
    } catch (error) {
      setResponse(`❌ 网络错误: ${error instanceof Error ? error.message : '未知错误'}\n\n请确保后端服务正在运行 (http://localhost:8080)`)
    } finally {
      setLoading(false)
    }
  }

  // 查询错误日志的方法
  const getErrorLogs = async () => {
    setLoading(true)
    setResponse('')
    
    try {
      const params = new URLSearchParams({
        project_id: 'demo-project',
        start_time: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(), // 24小时前
        end_time: new Date().toISOString()
      })

      const response = await fetch(`http://localhost:8080/api/error-logs?${params}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        }
      })

      if (response.ok) {
        const result = await response.json()
        setResponse(`查询到 ${result.length} 条错误日志: ${JSON.stringify(result, null, 2)}`)
      } else {
        setResponse(`请求失败: ${response.status} ${response.statusText}`)
      }
    } catch (error) {
      setResponse(`请求出错: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <div className="app-container">
        <h1>Spectra 日志收集演示</h1>
        
        <div className="demo-section">
          <h2>API 接口测试</h2>
          
          <div className="button-group">
            <button onClick={recordErrorLog} disabled={loading}>
              {loading ? '处理中...' : '记录错误日志'}
            </button>
            
            <button onClick={getErrorLogs} disabled={loading}>
              {loading ? '查询中...' : '查询错误日志'}
            </button>
            
            <button onClick={() => setCount(count + 1)}>
              增加计数器: {count}
            </button>
          </div>
          
          {response && (
            <div className="response-section">
              <h3>响应结果:</h3>
              <pre className="response-content">{response}</pre>
            </div>
          )}
        </div>
      </div>

    </>
  )
}

export default App
