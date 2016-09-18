import React     from 'react'

import Get       from './get.jsx'
import Delete    from './delete.jsx'


export default class GetDelParent extends React.Component {
  constructor(props) {
    super(props)
    //this.state = {
    //    ids: []
    //}

    //this.getBtnEvt = this.getBtnEvt.bind(this)
    //this.delBtnEvt = this.delBtnEvt.bind(this)
  }

  /*
  componentWillMount() {
    //Only once
    console.log("[GetDelParent]:componentWillMount()")

    let newIDs = []
    for (let user of this.props.users){
      newIDs.push(user.id)
    }

    this.setState({
      ids: newIDs
    })

  }*/

  /*
  componentWillReceiveProps(nextProps) {
    //最初は実行されない。データセットのイベント時のみ
    console.log("[GetDelParent]:componentWillReceiveProps()")
    //this.propsだとまだデータは更新されない
    //this.props.value -> old
    //nextProps.value -> new one
    //console.log("  this.props.users:", this.props.users)
    //console.log("  nextProps.users:", nextProps.users)

    let newIDs = []
    for (let user of nextProps.users){
      newIDs.push(user.id)
    }

    this.setState({
      ids: newIDs
    })
  }*/

  render() {
    return (
      <div className='row'>
        <Get btn={this.props.btnGet} ids={this.props.ids} />
        <Delete btn={this.props.btnDel} ids={this.props.ids} />
      </div>
    )
  }
}

GetDelParent.propTypes = { ids: React.PropTypes.array }
GetDelParent.defaultProps = { ids: [] }



/*
  var getBtn = document.getElementById("getBtn");
  var delBtn = document.getElementById("delBtn");
  getBtn.addEventListener("click", getUserList, false);
  delBtn.addEventListener("click", delUser, false);

  //get select text
  function getSelectText(id){
    // Get Slect box text
    var slctIds = document.getElementById(id);
    var idx = slctIds.selectedIndex;
    //var value = slctIds.options[idx].value;
    var text  = slctIds.options[idx].text;

    return text;
  }

  //
  function getUserList(evt){
    //
    var url;
    var sendData;
    var id = getSelectText("slctIds");
    if (id == "All"){
        url = "/api/users";
    }else{
        url = "/api/users/" + id;
    }

    //create data

    //send
    hy.sendAjax(url, "get", "form", "")
  }

  function delUser(evt){
    //
    var id = getSelectText("slctDelIds");
    var url = "/api/users/" + id;

    //send
    hy.sendAjax(url, "delete", "form", "")
  }
*/
