import React from 'react'

export default class Timer extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
    this.state = {
        second: 0
    }
    //this.ticker = this.ticker.bind(this)
  }

  ticker() {
    //console.log("ticker()")
    this.setState({
        second: this.state.second + 1
    })
  }

  componentDidMount() {
    //Only once
    //console.log("componentDidMount()")
    this.timer = setInterval(() => this.ticker(), 1000)
  }

  componentWillUnmount() {
    //console.log("componentWillUnmount()")
    clearInterval(this.timer)
  }

  render() {
    //let aaa = 100
    //let bbb = 200  
    return (
      <div>
        <h4>Timer</h4>
        <div>secondsElapsed: {this.state.second}</div>
        <hr/>
      </div>
    )
  }
}
