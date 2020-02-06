/* eslint-disable consistent-return */
import React,
{useState,
  useEffect}
  from 'react';
import './filter.css';
import {Checkbox,
  Tooltip,
  Button,
} from '@patternfly/react-core/dist/js/components';
import {API_URL} from '../../constants';
import {InfoCircleIcon} from '@patternfly/react-icons';
import store from '../redux/store';
import {FETCH_TASK_SUCCESS} from '../redux/Actions/TaskActionType';

const Filter: React.FC = (props:any) => {
  const [, setTags] = useState();
  const [status, setStatus]=useState();
  const [max, setMax]=useState(18);
  const [, setX]=useState(1); //
  const [clear, setClear]=useState(' ');
  const [show, setShow]=useState('Show all');
  const tagitem :any =[{id: '1000', value: 'task', isChecked: false},
    {id: '1001', value: 'pipeline', isChecked: false},
    {id: '1002', value: 'verified', isChecked: false}];
  const taglist=(data:any) =>{
    data.map((it:any ) =>
      tagitem.push({id: String(it.id), value: it.name, isChecked: false}));
    return tagitem;
  };
  useEffect(() => {
    // const fetchData = async () => {
    fetch(`${API_URL}/tags`)
        .then((res) => res.json())
        .then((data) => setTags(taglist(data)));
    setStatus(tagitem);
    // eslint-disable-next-line
  }, []);
  // / dummy function for refreshing component
  const ddff=() =>{
    setX((x) => x+1);
  };
  //  formation of filter url for fetching task and pipelines
  const filterurl=(e:any) =>{
    let typeurl='all';
    let verifiedurl='all';
    let tagsurl = '';

    status.forEach((it:any) => {
      if (it.id === e.target.id) {
        it.isChecked=e.target.checked;
      }
    },
    );
    setStatus(status);
    ddff();
    const temptype = status.slice(0, 2);
    const tempverified=status.slice(2, 3);
    const temptags=status.slice(3);
    if (tempverified[0].isChecked === true) {
      verifiedurl='true';
    }
    if (temptype[0].isChecked === true) {
      typeurl='task';
    }
    if (temptype[1].isChecked === true) {
      typeurl ='pipeline';
    }
    if ((temptype[0].isChecked === true) &&
       (temptype[1].isChecked === true)) {
      typeurl = 'all';
    }
    temptags.forEach((it:any) =>{
      if (it.isChecked === true) {
        tagsurl+=it.value+'|';
      }
      // if (it.isChecked === true) {
      //   setClear('ClearAll');
      // }
    });
    let flag:any= false;
    status.forEach((it:any) =>{
      if (it.isChecked === true) {
        flag=true;
      }
    });
    if (flag === true) {
      setClear('ClearAll');
    } else {
      setClear(' ');
    }


    fetch(`${API_URL}/resources/${typeurl}/${verifiedurl}?tags=${tagsurl} `)
        .then((resp) => resp.json())
        .then((data) =>{
          store.dispatch(
              {
                type: FETCH_TASK_SUCCESS,
                payload: data,
              },
          );
        });
  };

  // /   for clearing all checkbox
  const cleartag=(e:any) =>{
    status.forEach((it:any) =>{
      if (it.isChecked === true) {
        it.isChecked=false;
      }
    });
    setStatus(status);
    ddff();
    filterurl(e);
    setClear('');
  };
  let showresource:any;
  if (status !== undefined) {
    const resource =status.slice(0, 2);
    showresource = resource.map((it: any) => (
      <div key = {it} style={{marginBottom: '0.5em'}}>
        <Checkbox
          onClick={filterurl}
          isChecked={it.isChecked}
          style={{width: '1.2em', height: '1.2em'}}
          label={it.value[0].toUpperCase()+it.value.slice(1)}
          value={it.value}
          name="type"
          id={it.id}
          aria-label="uncontrolled checkbox example"

        />
      </div>
    ));
  }
  let showverifiedtask:any;
  // jsx element for show verifiedtask
  if ( status !== undefined) {
    const verifiedtask=status.slice(2, 3);
    showverifiedtask = verifiedtask.map((it: any) => (
      <div key = {it} style={{marginBottom: '0.5em'}}>
        <Checkbox
          onClick={filterurl}
          isChecked={it.isChecked}
          style={{width: '1.2em', height: '1.2em'}}
          label={it.value[0].toUpperCase()+it.value.slice(1)}
          value={it.value}
          name="verified"
          id={it.id}
          aria-label="uncontrolled checkbox example"

        />
      </div>
    ));
  }
  // jsx element for show tags
  let showTags:any ='';
  if (status !== undefined) {
    const tempstatus=status.slice(3, max);
    tempstatus.sort((a:any, b:any) =>
  (a.value> b.value) ? 1 :
  ((b.value > a.value) ? -1 : 0));
    showTags=
     tempstatus.map((it: any) => (
       <div key = {it} style={{marginBottom: '0.5em'}}>
         <Checkbox
           onClick={filterurl}
           isChecked={it.isChecked}
           style={{width: '1.2em', height: '1.2em'}}
           label={it.value[0].toUpperCase()+it.value.slice(1)}
           value={it.value}
           name="Tags"
           id={it.id}
           aria-label="uncontrolled checkbox example"

         />
       </div>
     ));
  }
  // /  for showing more
  const seeMore=(e:any) =>{
    const x=max;

    if (x< status.length) {
      setMax(status.length);
      setShow('Show less');
    } else {
      setMax((x) => x-13);
      setShow('Show all');
    }
  };

  return (
    <div className="filter-size" >
      <h2 style={{marginBottom: '1em'}}>
        {' '}
        <Button component='a' variant='link'
          onClick={cleartag}> {clear} </Button>{'  '}
      </h2>
      <h2 style={{marginBottom: '1em'}}>
        {' '}
        <b>Types</b>{'  '}
      </h2>
      {showresource}
      <h2 style={{marginBottom: '1em', marginTop: '1em'}}>
        {' '}
        <b>Verified </b>{'  '}
        <Tooltip content={<div>
           Verified Task and Pipelines by Tekton Catalog</div>}>
          <InfoCircleIcon />
        </Tooltip>
      </h2>
      {showverifiedtask}
      <h2 style={{marginBottom: '1em', marginTop: '1em'}}><b> Tags </b></h2>
      {showTags}
      <Button component='a' variant='link'
        onClick={seeMore}>{show} </Button>

    </div>
  );
};
export default Filter;

