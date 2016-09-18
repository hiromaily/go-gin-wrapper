import React from 'react'

export default class Counter extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
    this.state = {
        count: 0
    }

    this.countUp = this.countUp.bind(this)
  }

  countUp() {
    console.log("countUp()")
    this.setState({
        count: this.state.count + 1
    })
  }

  render() {
    return (
      <div>
        <h4>Counter 1.2.3...</h4>
        <p>{this.state.count}</p>
        <button onClick={this.countUp}>count up</button>
        <hr/>
      </div>
    )
  }
}
