import React from 'react';
import HandleClick from './HandleClick';
import Countdown from 'react-countdown';
import string from 'prop-types';

class MyCountdown extends React.Component {
  //   // Uninitialized state will cause Child to error out
  //   const [whenAutoPowerOff, setItems] = useState();

  //   // Data does't start loading
  //   // until *after* Parent is mounted
  //   useEffect(() => {
  //     fetch('/api/v1/get-time-autopoweroff/')
  //       .then(res => res.json())
  //       .then(data => setItems(data));
  //   }, []);

  // console.log(whenAutoPowerOff);
  render() {

    const renderer = ({ hours, minutes, seconds, completed }) => {
      if (completed) {
        HandleClick("shutdown", (new Date()).toISOString());
      } else {
        return <span>{hours}:{minutes}:{seconds}</span>;
      }
    };

    const dateNow = new Date()
    const time = dateNow.setTime(dateNow.getTime() + props.hours * 60 * 60 * 1000);
    
    return <Countdown date={time} renderer={renderer} />
  }
}

MyCountdown.propTypes = {
  hours: string,
  minutes: string,
  seconds: string,
  completed: Boolean,
};

export default MyCountdown;