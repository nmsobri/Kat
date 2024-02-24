# Kat Programming Language ( Work In Progress )

First stage, the code will be evaluated at run time ( interpreted )  
Second stage would be compile it to custom byte code and then run it  
Hopefully it will run the following code  

```
let fmt = import("fmt")

struct User {
    name,
    age,
    job,

    fn info(self) {
        fmt.Print("Name is %s, age is %d, job is %s", self.name, self.age, self.job)
    }
}

fn main() {
    let arr = [1,2,3,4,5]
    let map = {name: "sobri", location: "penang"}

    let user = User{name:"Sobri", age:42, job:"Programmer"}
    user.info()

    say("hello world")

    let i =  1 * 2 + 1

    if i < 3 {
        fmt.Print("lower than 3\n")
    } else if i < 5 {
        fmt.Print("lower than 5\n")
    } else if i < 10 {
        fmt.Print("lower than 10\n")
    }

    for i > 0 {
        fmt.Print(i)
        i--
    }

    for let j = 0; j < 5; j++ {
        fmt.Print(j)
    }
}

fn say(greet) {
    fmt.Print(greet)
}
```
