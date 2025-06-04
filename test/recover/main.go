package main

func main(){
    defer println("defer 1")
    level1()
}



func level1(){
    defer println("defer 2")
    defer func(){
        if err := recover(); err != nil {
            println("recover in level1")
        }
    }()
    defer println("defer 3")
    level2()
}

func level2(){
    defer println("defer 4")
    panic("panic in level2")
}