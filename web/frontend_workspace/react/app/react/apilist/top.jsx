import React     from 'react'
import ReactDOM  from 'react-dom'
import $         from 'jquery'

import List          from './components/list/parent.jsx'
import JwtParent     from './components/jwt/parent.jsx'
import GetDelParent  from './components/getdel/parent.jsx'
import PutPostParent from './components/putpost/parent.jsx'

export default class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
        users: [],
        ids: []
    }

    this.getBtnEvt = this.getBtnEvt.bind(this)
    this.delBtnEvt = this.delBtnEvt.bind(this)
    this.putBtnEvt = this.putBtnEvt.bind(this)
    this.postBtnEvt = this.postBtnEvt.bind(this)
    this.getUserIDs = this.getUserIDs.bind(this)
  }

  //Only once before first render()
  componentWillMount() {
    console.log("[App]:componentWillMount()")
    this.getUserIDs()
  }

  callAjax(mode, passedURL, sendData) {
    console.log("[App]:callAjax, mode is ", mode)

    let that = this
    let url = passedURL
    let method = 'get'
    let contentType = "application/x-www-form-urlencoded"
    let token = "Bearer " + this.refs.jwt.state.code

    switch (mode) {
      case 'getid':
        break
      case 'getlist':
        break
      case 'delete':
        method = 'delete'
        break
      case 'put':
        method = 'put'
        break
      case 'post':
        method = 'post'
        break
      default:
        return
        break
    }

    $.ajax({
      url: encodeURI(url),
      type: method,
      beforeSend: function beforeSend(xhr) {
        xhr.setRequestHeader(hiromaily_header, hiromaily_key)
        xhr.setRequestHeader('Authorization', token)
      },
      //cache    : false,
      crossDomain: false,
      //contentType: contentType,
      dataType:    'json', //data type from server
      data:        sendData
    }).done(function (data, textStatus, jqXHR) {
      switch (mode) {
        case 'getid':
          //console.log(data.ids)
          that.setState({ids: data.ids})
          //swal("success!", "get ids!", "success")
          break
        case 'getlist':
          if (data.code == 0) {
            //console.log(data.users)
            //let newUsers = []
            //for (let user of data.users) {
            //  newUsers.push({id:user.id, firstName:user.firstName, lastName:user.lastName, email:user.email, updated:user.updated})
            //}
            that.setState({
              users: data.users
            })
            //swal("success!", "get user list!", "success")
          }else{
            swal("error!", "validation error was occurred!", "error")
          }
          break
        case 'delete':
          console.log(data.id)
          // get user ids again.
          that.getUserIDs()
          that.getUsers('All')
          swal("success!", "delete user!", "success")
          break
        case 'put':
          console.log(data.id)
          that.getUsers(data.id)
          swal("success!", "put user!", "success")
          break
        case 'post':
          console.log(data.id)
          that.getUserIDs()
          that.getUsers(data.id)
          swal("success!", "post user!", "success")
          break
        default:
          break
      }
    }).fail(function (jqXHR, textStatus, errorThrown) {
      console.error(url, textStatus, errorThrown.toString())
      swal("error!", "validation error was occurred!", "error")
    })
  }

  //Get all user IDs
  getUserIDs() {
    console.log("[App]:getUserIDs()")
    //JWT
    if(this.refs.jwt == undefined){
      return
    }else if (this.refs.jwt.state.code == ""){
      return
    }
    //console.log("jwt code:", this.refs.jwt.state.code)

    //call ajax
    let url = '/api/users/ids'
    this.callAjax('getid', url, '')
  }

  //Get users data as list
  getUsers(id) {
    console.log("[App]:getUsers()")

    if (this.refs.jwt.state.code == ""){
      swal("warning!", "jwt code is required.", "warning")
      return
    }

    let url = ''
    if (id == "All"){
        url = "/api/users"
    }else{
        url = "/api/users/id/" + id
    }

    //call ajax
    this.callAjax('getlist', url, '')

    /*
    var newUsers = []
    newUsers.push({id:1, firstName:'harry1', lastName:'yasu1', email:'aa@aa.jp', updated:'2016/5/21'})
    newUsers.push({id:2, firstName:'harry2', lastName:'yasu2', email:'aa@aa.jp', updated:'2016/5/21'})
    newUsers.push({id:3, firstName:'harry3', lastName:'yasu3', email:'aa@aa.jp', updated:'2016/5/21'})
    newUsers.push({id:4, firstName:'harry4', lastName:'yasu4', email:'aa@aa.jp', updated:'2016/5/21'})
    newUsers.push({id:5, firstName:'harry5', lastName:'yasu5', email:'aa@aa.jp', updated:'2016/5/21'})
    newUsers.push({id:6, firstName:'harry6', lastName:'yasu6', email:'aa@aa.jp', updated:'2016/5/21'})

    this.setState({
      users: newUsers
    })
    */
  }

  //Delete users data from id
  deleteUser(id) {
    console.log("[App]:deleteUser()")

    if (this.refs.jwt.state.code == ""){
      swal("warning!", "jwt code is required.", "warning")
      return
    }

    let url = "/api/users/id/" + id

    //call ajax
    this.callAjax('delete', url, '')
  }

  //Put user data
  putUser(data) {
    console.log("[App]:putUser()")

    if (this.refs.jwt.state.code == ""){
      swal("warning!", "jwt code is required.", "warning")
      return
    }

    let url = "/api/users/id/"
    let sendData = new Object()

    if (data.id != "" && !isNaN(data.id)){
      url = url + data.id
    }else{
      swal("warning!", "user id is invalid.", "warning")
      return
    }


    //create data
    if (data.firstName != ""){
      sendData.firstName = data.firstName
    }
    if (data.lastName != ""){
      sendData.lastName = data.lastName
    }
    if (data.email != ""){
      sendData.email = data.email
    }
    if (data.password != ""){
      sendData.password = data.password
    }

    //call ajax
    this.callAjax('put', url, sendData)
  }

  //Post user data
  postUser(data) {
    console.log("[App]:postUser()")

    if (this.refs.jwt.state.code == ""){
      swal("warning!", "jwt code is required.", "warning")
      return
    }

    let url = "/api/users"
    let sendData = new Object()

    //create data
    sendData.firstName = data.firstName
    sendData.lastName = data.lastName
    sendData.email = data.email
    sendData.password = data.password

    //call ajax
    this.callAjax('post', url, sendData)
  }

  //Click get btn
  getBtnEvt(id) {
    console.log("[App]:getBtnEvt()")
    //console.log(id)

    //1.get user list by id
    this.getUsers(id)

    //[!!!]set.State is not updated until calling render()
    //console.log("[App]:getBtnEvt(), this.state", this.state)
  }

  //Click delete btn
  delBtnEvt(id) {
    console.log("[App]:delBtnEvt()")
    //console.log(id)

    //1.delete user by id
    this.deleteUser(id)
  }

  //Click put btn
  putBtnEvt(data) {
    console.log("[App]:putBtnEvt()")
    //console.log(data)

    //put user
    this.putUser(data)
  }

  //Click post btn
  postBtnEvt(data) {
    console.log("[App]:postBtnEvt()")
    //console.log(data)

    //post user
    this.postUser(data)
  }


  render() {
    //state was updated at this time.
    //console.log("[App]:render() this.state.users:", this.state.users)

    return (
      <div>
        <List users={this.state.users} />
        <JwtParent users={this.state.users} funcGetUserIDs={this.getUserIDs} ref="jwt" />
        <GetDelParent btnGet={this.getBtnEvt} btnDel={this.delBtnEvt} ids={this.state.ids} ref="getdel" />
        <PutPostParent btnPut={this.putBtnEvt} btnPost={this.postBtnEvt} />
      </div>
    )
  }
}

ReactDOM.render(
  <App />,
  document.getElementById('app')
)
