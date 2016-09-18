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
  }

  //Only once before first render()
  componentWillMount() {
    console.log("[App]:componentWillMount()")
    this.getUserIDs()
  }

  //Get all user IDs
  //TODO:get user ids by Ajax
  getUserIDs() {
    console.log("[App]:getUserIDs()")
    console.log(hiromaily_header)
    console.log(hiromaily_key)
    //JWT
    if(this.refs.jwt != undefined){
      console.log(this.refs.jwt.state.code)
    }

    let url = '/json/userIDs.json'

    //Only this API can Access without jwt
    $.ajax({
      url: url,
      dataType: 'json',
      cache: false,
      success: (data) => {
        this.setState({ids: data.ids})
      },
      error: (xhr, status, err) => {
        console.error(url, status, err.toString())
      }
    })
    //this.setState({
    //  ids: [1,2,3,4,5,6,7,8,9,10]
    //})
  }

  //This may be not necessary
  getUserIDsFromUsers(users){
    let newIDs = []
    for (let user of users){
      newIDs.push(user.id)
    }

    this.setState({
      ids: newIDs
    })
  }

  //Get users data as list
  getUsers() {
    //TODO:get users by Ajax
    console.log("[App]:getUsers()")

    var newUsers = []
    newUsers.push({id:1, firstName:'harry1', lastName:'yasu1', email:'aa@aa.jp', updated:'2016/5/21'})
    newUsers.push({id:2, firstName:'harry2', lastName:'yasu2', email:'aa@aa.jp', updated:'2016/5/21'})
    newUsers.push({id:3, firstName:'harry3', lastName:'yasu3', email:'aa@aa.jp', updated:'2016/5/21'})
    newUsers.push({id:4, firstName:'harry4', lastName:'yasu4', email:'aa@aa.jp', updated:'2016/5/21'})
    newUsers.push({id:5, firstName:'harry5', lastName:'yasu5', email:'aa@aa.jp', updated:'2016/5/21'})
    newUsers.push({id:6, firstName:'harry6', lastName:'yasu6', email:'aa@aa.jp', updated:'2016/5/21'})

    this.setState({
      users: newUsers
    });
  }

  //Click get btn
  getBtnEvt(data) {
    console.log("[App]:getBtnEvt()")
    console.log(data)
    //TODO: get userlist by Ajax
    //TODO: set response user list
    //たとえ1レコードでも全ユーザーIDを取得しないといけない。

    //1.指定したIDの情報取得 -> users
    this.getUsers()

    //2.userid全リストの取得 -> ids (不要)
    //this.getUserIDs()

    //set.Stateはrenderが呼ばれるまで反映されない。
    //console.log("[App]:getBtnEvt(), this.state", this.state)
  }

  //Click delete btn
  delBtnEvt(data) {
    console.log("[App]:delBtnEvt()")
    console.log(data)
    //TODO: delete user by Ajax
    //TODO: change user id options and user list

    //削除後、userlistを取得し直す
    //2.指定したIDの情報取得 -> users
    this.getUsers()

    //削除後、useridを取得し直す
    //3.userid全リストの取得 -> ids
    this.getUserIDs()

  }

  //Click put btn
  putBtnEvt(data) {
    console.log("[App]:putBtnEvt()")
    console.log(data)
    //TODO: update user by Ajax
    //TODO: change user id options and user list



    //更新後、userlistを取得し直す
    //2.指定したIDの情報取得 -> users
    this.getUsers()

  }

  //Click post btn
  postBtnEvt(data) {
    console.log("[App]:postBtnEvt()")
    console.log(data)
    //TODO: insert user by Ajax
    //TODO: change user id options and user list

    //登録後、userlistを取得し直す
    //2.指定したIDの情報取得 -> users
    this.getUsers()

    //後、useridを取得し直す
    //3.userid全リストの取得 -> ids
    this.getUserIDs()
  }


  render() {
    //このタイミングではstateは更新される
    //console.log("[App]:render() this.state.users:", this.state.users)

    return (
      <div>
        <List users={this.state.users} />
        <JwtParent users={this.state.users} ref="jwt" />
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
