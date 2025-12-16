package logger

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"

	pkgerrors "github.com/pkg/errors"
)

// ========================================
// 1. インターフェース定義
// ========================================

// stackTracer pkg/errorsのStackTrace interface
// pkg/errorsのエラーがこのインターフェースを実装している
type stackTracer interface {
	StackTrace() pkgerrors.StackTrace
}

// ========================================
// 2. メイン関数：ErrorAttr
// ========================================

// ErrorWithStackTrace エラーをslog属性に変換（stacktrace付き）
//
// 使い方:
//
//	logger.ErrorContext(ctx, "処理失敗", logger.ErrorWithStackTrace(err))
//
// 出力例:
//
//	{
//	  "error": {
//	    "message": "ユーザー登録失敗: DB接続エラー",
//	    "file": "/app/repository/user.go",
//	    "line": 42,
//	    "function": "repository.Insert",
//	    "stack_trace": "repository.Insert\n\t/app/repository/user.go:42\n..."
//	  }
//	}
func ErrorWithStackTrace(err error) slog.Attr {
	// nilチェック
	if err == nil {
		return slog.Attr{}
	}

	// 基本情報：エラーメッセージ
	attrs := []any{
		slog.String("message", err.Error()),
	}

	// pkg/errorsのstacktraceを取得
	st := getStackTrace(err)
	if st != nil {
		frames := st.StackTrace()
		if len(frames) > 0 {
			// エラー発生箇所の情報を追加
			attrs = append(attrs, extractErrorLocation(frames[0])...)

			// 全スタックトレースを追加
			attrs = append(attrs,
				slog.String("stack_trace", formatStackTrace(frames)),
			)
		}
	}

	// "error"グループにまとめて返す
	return slog.Group("error", attrs...)
}

// ========================================
// 3. スタックトレースの取得
// ========================================

// getStackTrace エラーからスタックトレースを取得
//
// エラーが複数回ラップされている場合：
//
//	err1 := pkgerrors.New("元のエラー")           // stacktrace: [A, B, C]
//	err2 := fmt.Errorf("ラップ1: %w", err1)      // stacktraceなし
//	err3 := pkgerrors.Wrap(err2, "ラップ2")       // stacktrace: [X, Y, Z]
//
// このとき、最初のstacktrace [A, B, C] を返す（最も情報量が多い）
func getStackTrace(err error) stackTracer {
	var firstStackTrace stackTracer

	// エラーチェーンを辿る
	for {
		// このエラーがstacktraceを持っているか確認
		if st, ok := err.(stackTracer); ok {
			// 最初に見つけたstacktraceを保存
			// （後でラップされたものより、元のものの方が詳細）
			firstStackTrace = st
		}

		// 次のエラーに進む
		err = errors.Unwrap(err)
		if err == nil {
			// これ以上エラーがない
			return firstStackTrace
		}
	}
}

// ========================================
// 4. スタックトレース情報の抽出
// ========================================

// extractErrorLocation フレームからエラー発生箇所の情報を抽出
//
// pkg/errorsのFrameからファイル名、行番号、関数名を取得する
func extractErrorLocation(frame pkgerrors.Frame) []any {
	// フレームからプログラムカウンタ（PC）を取得
	// PCは「プログラムのどこを実行しているか」を示す値
	pc := uintptr(frame) - 1

	// PCから関数情報を取得
	fn := runtime.FuncForPC(pc)

	file := "unknown"
	line := 0
	funcName := "unknown"

	if fn != nil {
		// ファイル名と行番号を取得
		file, line = fn.FileLine(pc)
		// 関数名を取得
		funcName = fn.Name()
	}

	return []any{
		slog.String("file", file),
		slog.Int("line", line),
		slog.String("function", funcName),
	}
}

// formatStackTrace スタックトレースを文字列にフォーマット
//
// pkg/errorsの%+vフォーマットを使用
// 出力例:
//
//	repository.Insert
//	    /app/repository/user.go:42
//	service.CreateUser
//	    /app/service/user.go:25
//	handler.HandleRequest
//	    /app/handler/user.go:15
func formatStackTrace(frames pkgerrors.StackTrace) string {
	return fmt.Sprintf("%+v", frames)
}

