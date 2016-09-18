import React from 'react'

// Here Increment and Decrement are being represented as stateless functional components. Note that props is passed as an argument & 'this' is not required. Also we don't need to implement render(), just return some valid JSX.
const Increment = (props) => {
      return <button onClick={props.increment}>Higher</button>
}

const Decrement = (props) => {
      return <button onClick={props.decrement}>Lower</button>
}

export default class Counter3 extends React.Component {
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

  
