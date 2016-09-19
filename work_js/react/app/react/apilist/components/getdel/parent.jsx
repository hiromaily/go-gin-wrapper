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
