import React     from 'react'
import ReactDOM  from 'react-dom'

import Header       from './components/layouts/header.jsx'
import Footer       from './components/layouts/footer.jsx'
import MemoBox      from './components/memos/memo_box.jsx'
import MemoBox2     from './components/memos/memo_box2.jsx'
import Timer        from './components/timer/timer.jsx'
import Counter      from './components/counter/counter.jsx'
import Counter2     from './components/counter/counter2.jsx'
import Counter3     from './components/counter/counter3.jsx'
import NewsBox      from './components/news/news.jsx'
import ColorBox     from './components/colorbox/colorbox.jsx'
import Communicate  from './components/communicate/communicate.jsx'
import User         from './components/context/context.jsx'
import User2        from './components/context/context2.jsx'

var memoData = [
  {id:1, language: "Golang", rank: 1},
  {id:2, language: "Python", rank: 2},
  {id:3, language: "Javascript", rank: 3}
]


export default class App extends React.Component {
  render() {
    return (
      <div className='contents'>
        <Header />
        <br/>  
        <Timer />
        <br/>  
        <MemoBox />
        <br/>  
        <NewsBox url="/json/news" pollInterval={10000} />
        <br/>
        <ColorBox />
        <br/>
        <Communicate />
        <br/>
        <Counter />
        <br/>
        <Counter2 />
        <br/>
        <Counter3 />
        <br/>
        <User />
        <br/>
        <User2 />
        <br/>
        <MemoBox2 data={memoData} />
        <br/>  
        <Footer />
      </div>
    )
  }
}

ReactDOM.render(
  <App />,
  document.getElementById('root')
)
