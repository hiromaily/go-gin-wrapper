import React from 'react'

export default class ColorBox extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
    this.state = {
        color: '#000000'
    }
  }

  changeColor() {
    console.log("[ColorBox] changeColor()")
    this.setState({
      color: "#000000".replace(/0/g,function(){return (~~(Math.random()*16)).toString(16);})
    })
  }

  componentDidMount() {
    console.log("[ColorBox] componentDidMount()")
    //setInterval(this.changeColor.bind(this), 2500)
    this.timer = setInterval(() => this.changeColor(), 2500)
  }
  componentWillUnmount() {
    clearInterval(this.timer)
  }

  render() {
    return(
      <div>
        <div className="block" style={{backgroundColor: this.state.color}}></div>
        <hr/>
      </div>
    )
  }
}
