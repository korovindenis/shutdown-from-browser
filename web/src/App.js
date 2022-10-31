import React, { Component, useState, useEffect } from 'react';
import homer from './homer.webp';
import './App.css';
import HandleClick from './components/HandleClick';
import MyButton from './components/MyButton';
import MyCountdown from './components/MyCountdown';
import { makeStyles } from '@material-ui/core/styles';
import Slider from '@material-ui/core/Slider';
import Typography from '@material-ui/core/Typography';
import { MDBFooter } from 'mdb-react-ui-kit';

// Button Style
const useStyles = makeStyles((theme) => ({
  cssProp: {
    margin: theme.spacing(1),
    fontSize: 20,
    color: '#FFF',
    width: '233px',
    height: '64px'
  },
  colorGreen: {
    backgroundColor: '#52b202',
  },
  colorRed: {
    backgroundColor: '#ff1744',
  },
  leftText: {
    textAlign: "left"
  },
  autoPowerOff: {
    margin: '15px'
  }
}));

function App() {
  const timeFormatted = (t) => { return t}
  const classes = useStyles();
  const buttons = [
    {
      text: "reboot",
      css: `${classes.cssProp} ${classes.colorGreen}`
    },
    {
      text: "shutdown",
      css: `${classes.cssProp} ${classes.colorRed}`
    }
  ]

  const compareTime = (dt) => {
    const res = Math.abs(new Date(dt) - new Date(new Date().toUTCString())) / 36e5 || 0
    return res;
  }
  const [whenAutoPowerOff, setItems] = useState();
  useEffect(() => {
    fetch('/api/v1/get-time-autopoweroff/')
      .then(res => res.json())
      .then(data => {setItems(data);sliderChangeValue(null,compareTime(data.time))});
  }, []);

  const [autoPowerOff, setautoPowerOff] = useState("is disabled");
  const sliderChangeValue = (e, value) => {
    HandleClick(value == 0 ? "" : "shutdown", new Date(new Date().getTime() + value * 60 * 60 * 1000).toISOString().replace(/:[0-9]{2}\.[0-9]{3}Z/,":00.000Z"))

    const _myCountdown = <MyCountdown hours={value} />
    if (value > 0) {
      setautoPowerOff(_myCountdown)
    } else {
      setautoPowerOff("is disable")
    }
  };

  return (
    <div className="App">
      <div className="App-main">
        <img src={homer} className="App-logo" alt="logo" />
        <div>
          {buttons.map(button => (
            <MyButton key={button.text} text={button.text} css={button.css} />
          ))}
          <div className={classes.autoPowerOff}>
            {!whenAutoPowerOff ? <p>Loading...</p> :
              <div>
                <Typography className={classes.leftText}>
                  {/* Auto-PowerOff {cTime > 0 ? <MyCountdown hours={cTime} /> : autoPowerOff} */}
                  Auto-PowerOff {autoPowerOff}
                </Typography>
                <Slider
                  defaultValue={compareTime(whenAutoPowerOff.time)}
                  aria-labelledby="discrete-slider"
                  valueLabelDisplay="auto"
                  onChange={sliderChangeValue}
                  step={1}
                  marks
                  min={0}
                  max={24}
                />
              </div>
            }
          </div>
        </div>
        <MDBFooter bgColor='light' className='text-center text-lg-start text-muted App-footer'>
          <div className='text-center p-3'>
            <a className='text-white' target="_blank" rel="noopener noreferrer" href='https://github.com/korovindenis/shutdown-from-browser'>
              github.com/korovindenis
            </a>
          </div>
        </MDBFooter>
      </div>
    </div>
  );
}
export default App;
