import React from 'react'
import $ from 'jquery'
import NewsList   from './newslist.jsx'

export default class NewsBox extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      data: [],
      issueNum: 1
    }
    //this.issueNum = 1

  }
  //defaultProps can't be used in React.Component

  // Ajax
  loadCommentsFromServer() {
    console.log("[NewsBox] loadCommentsFromServer()")
    //this stands for parent?
    $.ajax({
      url: this.props.url + `${this.state.issueNum}.json`,
      dataType: 'json',
      cache: false,
      success: (data) => { 
        this.setState({data: data.news})
      },
      error: (xhr, status, err) => {
        console.error(this.props.url, status, err.toString())
      }
    })

    if (this.state.issueNum+1 == 4){
      this.setState({ issueNum: 1 })      
    }else{
      this.setState({ issueNum: this.state.issueNum + 1 })      
    }

    //this.issueNum++
    //if (this.issueNum == 4) {
    //  this.issueNum = 1
    //}
  }

  //
  componentDidMount() {
    console.log("[NewsBox] componentDidMount()")
    this.loadCommentsFromServer()
    setInterval(this.loadCommentsFromServer.bind(this), this.props.pollInterval)
  }

  //
  render() {
    //loop for multiple data
    //for(let val of this.state.data) {
    //  console.log(val)
    //} 
    //*
    let newsNodes = this.state.data.map(function(news) {
      return (
        //<NewsList data={this.state.data} />
        <NewsList key={news.id} date={news.date} title={news.title} author={news.author}>
          {news.text}
        </NewsList>
      )
    })
    //*/
    return(
      <div className='newsBox'>
        <h4>News</h4>
        {newsNodes}
        <hr/>
      </div>
    )
  }
}
//defaultProps can't be used in React.Component
//NewsBox.defaultProps = { issueNum: 1 }