import { useState } from 'react'
import './App.css'

function App() {
  const [count, setCount] = useState(0)
  const [response, setResponse] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)

  const apiBase = 'http://localhost:8080/api'
  const nowISO = () => new Date().toISOString()
  const timeRange = () => ({
    start_time: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
    end_time: new Date().toISOString()
  })
  const baseLog = () => ({
    timestamp: nowISO(),
    project_id: 'demo-project',
    session_id: 'session-123',
    trace_id: `trace-${Date.now()}`,
    user_id: 'user-789',
    url: window.location.href,
    referrer: document.referrer,
  })

  // 记录错误日志的方法
  const recordErrorLog = async () => {
    setLoading(true)
    setResponse('')
    
    try {
      const errorData = {
        ...baseLog(),
        type: 'javascript_error',
        name: 'DemoError',
        message: '这是一个演示错误日志',
        extra: {
          browser: navigator.userAgent,
          count: count
        }
      }

      const response = await fetch(`${apiBase}/error-logs`, {
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
        ...timeRange()
      })

      const response = await fetch(`${apiBase}/error-logs?${params}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        }
      })

      if (response.ok) {
        const result = await response.json()
        if (result.length === 0) {
          setResponse(`ℹ️ 查询完成，但没有找到错误日志记录。\n\n查询参数:\n项目ID: demo-project\n时间范围: 过去24小时`)
        } else {
          setResponse(`✅ 查询成功！\n\n查询到 ${result.length} 条错误日志:\n${JSON.stringify(result, null, 2)}`)
        }
      } else {
        const errorText = await response.text()
        setResponse(`❌ 查询失败: ${response.status} ${response.statusText}\n\n错误详情:\n${errorText}`)
      }
    } catch (error) {
      setResponse(`❌ 网络错误: ${error instanceof Error ? error.message : '未知错误'}\n\n请确保后端服务正在运行 (http://localhost:8080)`)
    } finally {
      setLoading(false)
    }
  }

  // 记录性能指标
  const recordPerformanceMetric = async () => {
    setLoading(true)
    setResponse('')
    try {
      const data = {
        ...baseLog(),
        type: 'performance_metric',
        name: 'DemoMetric',
        value: Math.round(performance.now()),
        extra: {
          browser: navigator.userAgent,
          metric: 'perf'
        }
      }
      const res = await fetch(`${apiBase}/performance-metrics`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
      })
      if (res.ok) {
        const result = await res.json()
        setResponse(`✅ 成功记录性能指标！\n\n响应数据:\n${JSON.stringify(result, null, 2)}`)
      } else {
        setResponse(`❌ 请求失败: ${res.status} ${res.statusText}\n\n错误详情:\n${await res.text()}`)
      }
    } catch (error) {
      setResponse(`❌ 网络错误: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      setLoading(false)
    }
  }

  // 查询性能指标
  const getPerformanceMetrics = async () => {
    setLoading(true)
    setResponse('')
    try {
      const params = new URLSearchParams({ project_id: 'demo-project', ...timeRange() })
      const res = await fetch(`${apiBase}/performance-metrics?${params}`, { method: 'GET', headers: { 'Content-Type': 'application/json' } })
      if (res.ok) {
        const result = await res.json()
        setResponse(`✅ 查询性能指标成功！\n\n${JSON.stringify(result, null, 2)}`)
      } else {
        setResponse(`❌ 查询失败: ${res.status} ${res.statusText}\n\n错误详情:\n${await res.text()}`)
      }
    } catch (error) {
      setResponse(`❌ 网络错误: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      setLoading(false)
    }
  }

  // 记录用户行为
  const recordUserAction = async () => {
    setLoading(true)
    setResponse('')
    try {
      const data = {
        ...baseLog(),
        type: 'user_action',
        name: 'ButtonClick',
        message: '点击了测试按钮',
        method: 'CLICK',
        status: 200,
        value: 1,
        extra: { button: '记录用户行为', count }
      }
      const res = await fetch(`${apiBase}/user-actions`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
      })
      if (res.ok) {
        const result = await res.json()
        setResponse(`✅ 成功记录用户行为！\n\n响应数据:\n${JSON.stringify(result, null, 2)}`)
      } else {
        setResponse(`❌ 请求失败: ${res.status} ${res.statusText}\n\n错误详情:\n${await res.text()}`)
      }
    } catch (error) {
      setResponse(`❌ 网络错误: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      setLoading(false)
    }
  }

  // 查询用户行为
  const getUserActions = async () => {
    setLoading(true)
    setResponse('')
    try {
      const params = new URLSearchParams({ project_id: 'demo-project', ...timeRange() })
      const res = await fetch(`${apiBase}/user-actions?${params}`, { method: 'GET', headers: { 'Content-Type': 'application/json' } })
      if (res.ok) {
        const result = await res.json()
        setResponse(`✅ 查询用户行为成功！\n\n${JSON.stringify(result, null, 2)}`)
      } else {
        setResponse(`❌ 查询失败: ${res.status} ${res.statusText}\n\n错误详情:\n${await res.text()}`)
      }
    } catch (error) {
      setResponse(`❌ 网络错误: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      setLoading(false)
    }
  }

  // 记录自定义事件
  const recordCustomEvent = async () => {
    setLoading(true)
    setResponse('')
    try {
      const data = {
        ...baseLog(),
        type: 'custom_event',
        name: 'DemoEvent',
        message: '这是一个自定义事件',
        extra: { feature: 'demo', count }
      }
      const res = await fetch(`${apiBase}/custom-events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
      })
      if (res.ok) {
        const result = await res.json()
        setResponse(`✅ 成功记录自定义事件！\n\n响应数据:\n${JSON.stringify(result, null, 2)}`)
      } else {
        setResponse(`❌ 请求失败: ${res.status} ${res.statusText}\n\n错误详情:\n${await res.text()}`)
      }
    } catch (error) {
      setResponse(`❌ 网络错误: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      setLoading(false)
    }
  }

  // 查询自定义事件
  const getCustomEvents = async () => {
    setLoading(true)
    setResponse('')
    try {
      const params = new URLSearchParams({ project_id: 'demo-project', ...timeRange() })
      const res = await fetch(`${apiBase}/custom-events?${params}`, { method: 'GET', headers: { 'Content-Type': 'application/json' } })
      if (res.ok) {
        const result = await res.json()
        setResponse(`✅ 查询自定义事件成功！\n\n${JSON.stringify(result, null, 2)}`)
      } else {
        setResponse(`❌ 查询失败: ${res.status} ${res.statusText}\n\n错误详情:\n${await res.text()}`)
      }
    } catch (error) {
      setResponse(`❌ 网络错误: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      setLoading(false)
    }
  }

  // 记录页面停留
  const recordPageStay = async () => {
    setLoading(true)
    setResponse('')
    try {
      const data = {
        ...baseLog(),
        type: 'page_stay',
        name: document.title || 'DemoPage',
        value: Math.round(5 + Math.random() * 55), // 模拟停留秒数
        extra: { scrollTop: window.scrollY }
      }
      const res = await fetch(`${apiBase}/page-stays`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
      })
      if (res.ok) {
        const result = await res.json()
        setResponse(`✅ 成功记录页面停留！\n\n响应数据:\n${JSON.stringify(result, null, 2)}`)
      } else {
        setResponse(`❌ 请求失败: ${res.status} ${res.statusText}\n\n错误详情:\n${await res.text()}`)
      }
    } catch (error) {
      setResponse(`❌ 网络错误: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      setLoading(false)
    }
  }

  // 查询平均页面停留
  const getAveragePageStay = async () => {
    setLoading(true)
    setResponse('')
    try {
      const params = new URLSearchParams({ project_id: 'demo-project', ...timeRange() })
      const res = await fetch(`${apiBase}/page-stays/average?${params}`, { method: 'GET', headers: { 'Content-Type': 'application/json' } })
      if (res.ok) {
        const result = await res.json()
        setResponse(`✅ 查询平均页面停留成功！\n\n${JSON.stringify(result, null, 2)}`)
      } else {
        setResponse(`❌ 查询失败: ${res.status} ${res.statusText}\n\n错误详情:\n${await res.text()}`)
      }
    } catch (error) {
      setResponse(`❌ 网络错误: ${error instanceof Error ? error.message : '未知错误'}`)
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

            <button onClick={recordPerformanceMetric} disabled={loading}>
              {loading ? '处理中...' : '记录性能指标'}
            </button>

            <button onClick={getPerformanceMetrics} disabled={loading}>
              {loading ? '查询中...' : '查询性能指标'}
            </button>

            <button onClick={recordUserAction} disabled={loading}>
              {loading ? '处理中...' : '记录用户行为'}
            </button>

            <button onClick={getUserActions} disabled={loading}>
              {loading ? '查询中...' : '查询用户行为'}
            </button>

            <button onClick={recordCustomEvent} disabled={loading}>
              {loading ? '处理中...' : '记录自定义事件'}
            </button>

            <button onClick={getCustomEvents} disabled={loading}>
              {loading ? '查询中...' : '查询自定义事件'}
            </button>

            <button onClick={recordPageStay} disabled={loading}>
              {loading ? '处理中...' : '记录页面停留'}
            </button>

            <button onClick={getAveragePageStay} disabled={loading}>
              {loading ? '查询中...' : '查询平均页面停留'}
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
