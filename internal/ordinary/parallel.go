package ordinary

import "context"

// GroupFunc 调用函数
type GroupFunc func() error

// NewParallel 创建
func NewParallel(ctx context.Context, handles ...GroupFunc) error {
	errChan := make(chan error)
	doneChan := make(chan *struct{})

	for _, handle := range handles {
		currentHandle := handle

		go func() {
			if err := currentHandle(); err != nil {
				errChan <- err
			}

			doneChan <- nil
		}()
	}

	count := len(handles)

	for {
		select {
		case <-ctx.Done():
			{
				return ctx.Err()
			}

		case err := <-errChan:
			{
				return err
			}

		case <-doneChan:
			{
				count--
				if count <= 0 {
					return nil
				}
			}
		}
	}
}
