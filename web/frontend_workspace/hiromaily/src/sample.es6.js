'use strict'

var ES6 = 'ECMAScript2015'

// function log
let log = () => {
  console.log(ES6)
  console.log("ECMAScript2015")
}

// class Hy
class Hy {
    constructor(x, y) {
        //public
        this.x = x
        this.y = y
    }

    //public
    static distance(a, b) {
        const dx = a.x - b.x
        const dy = a.y - b.y

        return Math.sqrt(dx*dx + dy*dy)
    }
}

//main()
function main(){
  alert(11)

  log()

  //fmt.Printf的な
  let name = 'Harry'
  console.log(`Hello, ${name}`)

  //class
  const p1 = new Hy(5, 5)
  const p2 = new Hy(10, 10)

  console.log(Hy.distance(p1, p2))
  
}

main()
