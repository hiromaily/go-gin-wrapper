import React from 'react'
import Increment   from './increment.jsx'
import Decrement   from './decrement.jsx'

export default class Counter2 extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
    this.state = {
        currentCount: 0
    }
    //this.countUp = this.countUp.bind(this)
  }

  // Basic parent increment handler
  increment() {
    this.setState({
      currentCount: this.state.currentCount+1
    })
  }
  
  // Basic parent decrement handler
  decrement() {
    this.setState({
      currentCount: this.state.currentCount-1
    })
  }

  render() {
    //子のイベントを取得するためにbindしたfunctionを渡している。
    return (
      <div>
        <h2>Current Count: {this.state.currentCount}</h2>
        <Increment increment={this.increment.bind(this)} />
        <Decrement decrement={this.decrement.bind(this)} />
        <hr/>
      </div>
    )
  }
}

