# Kat Programming Language ( Work In Progress )

First stage, the code will be evaluated at run time ( interpreted )  
Second stage would be compile it to custom byte code and then run it  
Hopefully it will run the following code  

```go
const fmt = import("fmt")

struct User {
    name,
    age,
    job,
}

fn User.setAge(self, age) {
    self.age =  age
}

fn User.info(self) {
    fmt.Print("Name is %s, age is %d, job is %s", self.name, self.age, self.job)
}

fn say(greet) {
    fmt.Print(greet)
}

fn main() {
    let arr = [1,2,3,4,5]
    let foo = arr[0]
    let map = {name: "sobri", location: "penang"}
    let name = map["name"]

    let user = User{name:"Sobri", age:99, job:"Programmer"}
    user.setAge(99)
    user.info()

    say("hello world")

    let i =  1 * 2 + 1

    if 10 > 2 {
        fmt.Println("bigger than 2")
    } else {
        fmt.Println("smaller than 2")
    }

    if i < 3 {
        fmt.Print("lower than 3\n")
    } else if i < 5 {
        fmt.Print("lower than 5\n")
    } else if i < 10 {
        fmt.Print("lower than 10\n")
    }

    for i > 0 {
        fmt.Print(i)
        1 + i--
        1 + --i
		---3
		--1 + 2
		++5--
		--6--
    }

    for let j = 0; j < 5; j++ {
        fmt.Print(j)
    }
}
```
