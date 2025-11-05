// 错误日志表
CREATE TABLE error_logs
(
    timestamp   DateTime,
    project_id  String,
    session_id  String,
    trace_id    String,
    user_id     String,
    url         String,
    referrer    String,
    type        String,
    name        String,
    message     String,
    extra       JSON
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (project_id, timestamp);

// 性能指标表
CREATE TABLE performance_metrics
(
    timestamp   DateTime,
    project_id  String,
    session_id  String,
    trace_id    String,
    user_id     String,
    url         String,
    referrer    String,
    type        String,       -- performance
    name        String,       -- FCP / LCP / CLS / TTFB...
    value       Float64,      -- 指标数值（ms / 分数）
    extra       JSON
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (project_id, timestamp);

// 用户行为表
CREATE TABLE user_actions
(
    timestamp   DateTime,
    project_id  String,
    session_id  String,
    trace_id    String,
    user_id     String,
    url         String,
    referrer    String,
    type        String,      -- user
    name        String,      -- click / route / api_timing
    message     String,      -- 元素标识 / 路由信息
    method      String,      -- GET / POST（仅 api_timing）
    status      UInt16,      -- HTTP 状态码（仅 api_timing）
    value       Float64,     -- 接口耗时 / API duration
    extra       JSON
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (project_id, timestamp);

// 自定义事件表
CREATE TABLE custom_events
(
    timestamp   DateTime,
    project_id  String,
    session_id  String,
    trace_id    String,
    user_id     String,
    url         String,
    referrer    String,
    type        String,      -- custom
    name        String,      -- 自定义事件名
    message     String,      -- 固定为 custom_event
    extra       JSON         -- 自定义业务字段
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (project_id, timestamp);

// 页面停留时长表
CREATE TABLE page_stay
(
    timestamp   DateTime,
    project_id  String,
    session_id  String,
    trace_id    String,
    user_id     String,
    url         String,
    referrer    String,
    type        String,      -- page_stay
    name        String,      -- page_stay_time
    value       Float64,     -- 页面停留时长(ms)
    extra       JSON
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (project_id, timestamp);


