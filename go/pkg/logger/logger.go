package logger

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

var (
	defaultLogger *slog.Logger
	projectID     string
)

// Init ロガーの初期化
func Init(pid string) {
	projectID = pid
	SetDebug(true)
}

// SetDebug デバッグモードを設定
func SetDebug(debug bool) {
	var handler slog.Handler

	if debug {
		// 開発環境
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		// 本番環境: Cloud Logging形式（後で実装）
		// とりあえずJSON
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

func traceAttrs(ctx context.Context) []slog.Attr {
	if ctx == nil {
		return nil
	}

	sc := trace.SpanFromContext(ctx).SpanContext()
	if !sc.IsValid() {
		return nil
	}

	return []slog.Attr{
		slog.String("logging.googleapis.com/trace",
			fmt.Sprintf("projects/%s/traces/%s", projectID, sc.TraceID().String())),
		slog.String("logging.googleapis.com/spanId", sc.SpanID().String()),
		slog.Bool("logging.googleapis.com/trace_sampled", sc.IsSampled()),
	}
}

// ========================================
// 基本ログ関数
// サンプルコードのInfoWithContext等と同じ
// ========================================

// InfoContext Infoレベルのログ
func InfoContext(ctx context.Context, msg string, args ...any) {
	logWithTrace(ctx, slog.LevelInfo, msg, args...)
}

// WarnContext Warnレベルのログ
func WarnContext(ctx context.Context, msg string, args ...any) {
	logWithTrace(ctx, slog.LevelWarn, msg, args...)
}

// ErrorContext Errorレベルのログ
func ErrorContext(ctx context.Context, msg string, args ...any) {
	logWithTrace(ctx, slog.LevelError, msg, args...)
}

// DebugContext Debugレベルのログ
func DebugContext(ctx context.Context, msg string, args ...any) {
	logWithTrace(ctx, slog.LevelDebug, msg, args...)
}

// FatalContext Fatalレベルのログ（プログラム終了）
func FatalContext(ctx context.Context, msg string, args ...any) {
	logWithTrace(ctx, slog.LevelError, msg, args...)
	os.Exit(1)
}

// ========================================
// フォーマット付きログ関数
// サンプルコードのInfofWithContext等と同じ
// ========================================

// InfofContext Printf形式のInfoログ
func InfofContext(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	logWithTrace(ctx, slog.LevelInfo, msg)
}

// WarnfContext Printf形式のWarnログ
func WarnfContext(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	logWithTrace(ctx, slog.LevelWarn, msg)
}

// ErrorfContext Printf形式のErrorログ
func ErrorfContext(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	logWithTrace(ctx, slog.LevelError, msg)
}

// FatalfContext Printf形式のFatalログ
func FatalfContext(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	logWithTrace(ctx, slog.LevelError, msg)
	os.Exit(1)
}

// ========================================
// 内部ヘルパー
// ========================================

// logWithTrace トレース情報を含めてログ出力
// argsの中にerror型があれば、自動的にstacktraceを付ける
func logWithTrace(ctx context.Context, level slog.Level, msg string, args ...any) {
	// エラーを自動検出してstacktraceを付ける
	args = processErrorsInArgs(args)

	// トレース情報を取得
	if traceAttrs := traceAttrs(ctx); len(traceAttrs) > 0 {
		// トレース情報を引数の先頭に追加
		newArgs := make([]any, 0, len(args)+len(traceAttrs)*2)
		for _, attr := range traceAttrs {
			newArgs = append(newArgs, attr.Key, attr.Value.Any())
		}
		newArgs = append(newArgs, args...)
		args = newArgs
	}

	defaultLogger.Log(ctx, level, msg, args...)
}

// processErrorsInArgs argsの中のerror型を検出してstacktraceを付ける
func processErrorsInArgs(args []any) []any {
	newArgs := make([]any, 0, len(args))

	for i := 0; i < len(args); i += 2 {
		// 奇数個の場合、最後の要素をそのまま追加
		if i+1 >= len(args) {
			newArgs = append(newArgs, args[i])
			break
		}

		key := args[i]
		value := args[i+1]

		// error型かチェック
		if err, ok := value.(error); ok && err != nil {
			// ErrorWithStackTrace()で変換
			// 元のキーは削除して、"error"グループとして追加
			newArgs = append(newArgs, ErrorWithStackTrace(err))
		} else {
			// 通常のキー・バリュー
			newArgs = append(newArgs, key, value)
		}
	}

	return newArgs
}

// ========================================
// HTTP情報のログ
// サンプルコードのWithHTTPと同じ
// ========================================

// HTTPAttr HTTPリクエスト情報をslog属性として返す
func HTTPAttr(req *http.Request, statusCode int, duration time.Duration, responseSize int) slog.Attr {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}

	return slog.Group("httpRequest",
		slog.String("requestMethod", req.Method),
		slog.String("requestUrl", req.RequestURI),
		slog.Int("status", statusCode),
		slog.String("userAgent", req.UserAgent()),
		slog.String("remoteIp", host),
		slog.String("referer", req.Referer()),
		slog.String("protocol", req.Proto),
		slog.String("latency", fmt.Sprintf("%.9fs", duration.Seconds())),
		slog.Int("responseSize", responseSize),
	)
}

// ========================================
// 大きなログの分割
// サンプルコードのLargeInfoと同じ
// ========================================

const logSplitLength = 20000

type logSplit struct {
	UID         string `json:"uid"`
	Index       int    `json:"index"`
	TotalSplits int    `json:"totalSplits"`
}

// LargeInfo 大きなサイズのログを分割して出力する
//
// Cloud Loggingの制限により、一度に出力できるログのサイズに制限があるため
func LargeInfo(ctx context.Context, msg string) {
	rs := []rune(msg)
	rsLength := len(rs)
	logSplitUID := uuid.NewString()
	logSplitTotalSplits := int(math.Ceil(float64(rsLength) / float64(logSplitLength)))

	for i := 0; i < logSplitTotalSplits; i++ {
		start := logSplitLength * i
		end := start + logSplitLength
		if end > rsLength {
			end = rsLength
		}
		delimited := string(rs[start:end])

		split := logSplit{
			UID:         logSplitUID,
			Index:       i,
			TotalSplits: logSplitTotalSplits,
		}

		InfoContext(ctx, delimited,
			"type", "audit",
			"split", split,
		)
	}
}
