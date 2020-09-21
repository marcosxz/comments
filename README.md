# comments
go ast comment

```text
// 实现Finder接口
type printFinder struct{}
func (printFinder) Filter(os.FileInfo) bool {  ...  }
func (printFinder) Func(comments FuncComments) error {  ...  }
func (printFinder) Type(comments TypeComments) error {  ...  }
func (printFinder) Field(comments FieldComments) error {  ...  }

// 调用查找注释的方法
_ = comments.Find(dir, finder)
```
