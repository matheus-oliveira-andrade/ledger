package utils

type LoggerKey string
type CorrelationIdKey string

const CtxLoggerKey LoggerKey = "logger"
const CorrelationIdHeader CorrelationIdKey = "x-correlation-id"
