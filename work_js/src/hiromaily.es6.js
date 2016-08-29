'use strict'

// hy object
let hy = new Object
{
  //---------------------------------------------------------------------------
  //private
  //---------------------------------------------------------------------------
  let setToken = (token) => {
    var jwtCode = document.getElementById('jwtCode')
    jwtCode.value = token
  }

  let getTokenHeader = () => {
    var jwtCode = document.getElementById('jwtCode')
    return jwtCode.value
  }

  let updateUserList = (users) => {
    console.info('updateUserList()')
    let strHtml = ''
    users.forEach(user => {
      //console.log(users[i].id, users[i].firstName);
      strHtml += `<tr><td>${user.id}</td><td>${user.firstName}</td><td>${user.lastName}</td>
<td>${user.email}</td><td>*****</td><td>${user.update}</td>`
    })
    let userListBody = document.getElementById('userListBody')
    userListBody.innerHTML = strHtml
  }

  let square = (num) => {
    return num * num
  }

  //---------------------------------------------------------------------------
  //public
  //---------------------------------------------------------------------------
  hy.abc = () => {
    console.log(square(5))    
  }
}

// run
hy.abc()

