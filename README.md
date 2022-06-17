# go-zlog

package zap logger use with context

## Usage

```gol
func main() {
    //...
    ctx = zlog.With(ctx,  zap.String("request_id", requestID))
    CallFunc(ctx)
}

func CallFunc(ctx) {
    // ...
    zlog.Warnf(ctx, "some error occurred. err: %s", err) // zap.L().With(zap.String("request_id", requestID)).Sugar().Warnf("some error occurred. err: %s", err)
}
```
