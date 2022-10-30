# gosearch 

> 一个命令行下查询go mod的小工具

使用方式
1. 克隆代码
    ```
    git clone 
    ```
2. 安装
    ```
    go install cmd/gosearch.go
    ```
3. 使用
    ```
    gosearch [flags] keyword[[, ]keywords...]
    # examples
    gosearch gin
    gosearch gin,cobra
    gosearch gin cobra
    gosearch -s gin cobra
    gosearch -r gin cobra
    ```
