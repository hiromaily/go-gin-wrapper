import React from 'react'
import Message   from './message.jsx'

export default class MemoBox extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
    this.state = {
        memo: {
          title: 'nothing',
          day: 10
        }
    }
    //If passing this to another func, this point out parent one.
    this.chgTitle = this.chgTitle.bind(this)
    this.chgDay = this.chgDay.bind(this)
  }

  chgTitle(event) {
    console.log("chgTitle()")
    this.setState({
        memo: {
          title: event.target.value,
          day: this.state.memo.day
        }
    })
  }

  chgDay(event) {
    console.log("chgDay()")
    this.setState({
        memo: {
          title: this.state.memo.title,
          day: event.target.value
        }
    })
  }

  render() {
    return (
      <div className='memo_box'>
        <h4>Memo Box component</h4>
        <div>title:<input type="text" value={this.state.memo.title} onChange={this.chgTitle} /></div>
        <div>day:<input type="text" value={this.state.memo.day} onChange={this.chgDay} /></div>
        <Message title={this.state.memo.title} day={this.state.memo.day}>this is child element</Message>
        <hr/>
      </div>
    )
  }
}
