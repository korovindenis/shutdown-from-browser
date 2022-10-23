import Button from '@material-ui/core/Button';
import HandleClick from './HandleClick';

export default function MyButton(props) {
    return <Button 
      variant="contained" 
      className={props.css}
      onClick={e => HandleClick(props.text)}
      >
        {props.text.toUpperCase()}
    </Button>;
  }