/* eslint-disable consistent-return */
import React,
{
  useState,
  useEffect,
}
  from 'react';
import {
  Checkbox,
  Tooltip,
  Button,
} from '@patternfly/react-core/dist/js/components';
import {API_URL} from '../../constants';
import {InfoCircleIcon} from '@patternfly/react-icons';
import store from '../redux/store';
import {FETCH_TASK_SUCCESS} from '../redux/Actions/TaskActionType';
import './filter.css';

const Filter: React.FC = (props: any) => {
  const [categoryData, setCategoryData] = useState();
  const [status, setStatus] = useState({checklist: []});
  const [clear, setClear] = useState(' ');
  const filterItem: any = [{id: '1000', value: 'task', isChecked: false},
    {id: '1001', value: 'pipeline', isChecked: false},
    {id: '1002', value: 'verified', isChecked: false}];
  const [checkBoxStatus, setCheckBoxStatus] = useState(
      {
        task: false,
        pipeline: false,
        verified: false,
        C1: false,
        C2: false,
        C3: false,
        C4: false,
      },
  );

  //  function for adding categories to filteritem
  const addCategory = (categoryData: any) => {
    Object.keys(categoryData).map((categoryName: string, index: number) =>

      filterItem.push(
          {id: index.toString(), value: categoryName, isChecked: false},
      ));
    setStatus({checklist: filterItem});
    return categoryData;
  };
  useEffect(() => {
    fetch(`${API_URL}/categories`)
        .then((res) => res.json())
        .then((categoryData) => setCategoryData(addCategory(categoryData)));

    // eslint-disable-next-line
  }, []);

  // fetch api function
  const fetchApi = (typeurl: string, verifiedurl: string, categoryurl: any) => {
    fetch(`${API_URL}/resources/${typeurl}/${verifiedurl}?tags=${categoryurl} `)
        .then((resp) => resp.json())
        .then((data) => {
          const taskarr = data.sort((first:any, second: any) => {
            if (first.name > second.name) {
              return 1;
            } else {
              return -1;
            }
          });
          store.dispatch(
              {
                type: FETCH_TASK_SUCCESS,
                payload: taskarr,
              },
          );
        });
  };


  // formation of filter url  for calling filterAPi to
  //  fetching task and pipelines
  const filterApi = (event: any) => {
    let typeurl = 'all';
    let verifiedurl = 'all';
    let categoryurl = '';
    const target = event.target;
    // for handling isChecked parameter of checkbox
    setCheckBoxStatus({...checkBoxStatus, [target.value]: target.checked});
    status.checklist.forEach((it: any) => {
      if (it.id === event.target.id) {
        it.isChecked = event.target.checked;
      }
    },
    );

    const temptype = status.checklist.slice(0, 2);
    const tempverified = status.checklist.slice(2, 3);
    const tempcategory = status.checklist.slice(3);
    if (tempverified[0]['isChecked'] === true) {
      verifiedurl = 'true';
    }
    if (temptype[0]['isChecked'] === true) {
      typeurl = 'task';
    }
    if (temptype[1]['isChecked'] === true) {
      typeurl = 'pipeline';
    }
    if ((temptype[0]['isChecked'] === true) &&
      (temptype[1]['isChecked'] === true)) {
      typeurl = 'all';
    }
    tempcategory.map((item: any) => {
      if (item.isChecked === true) {
        categoryData[item.value].map((categoryitem: any) => {
          categoryurl = categoryurl + categoryitem + '|';

          return categoryurl;
        });
      }
      return tempcategory;
    });

    // for displaying clear filter options
    let flag: any = false;
    status.checklist.forEach((it: any) => {
      if (it.isChecked === true) {
        flag = true;
      }
    });
    if (flag === true) {
      setClear('ClearAll');
    } else {
      setClear(' ');
    }
    fetchApi(typeurl, verifiedurl, categoryurl);
  };

  //   function for clearing all checkbox
  const clearFilter = () => {
    setCheckBoxStatus(
        {
          task: false,
          pipeline: false,
          verified: false,
          C1: false,
          C2: false,
          C3: false,
          C4: false,
        });
    status.checklist.map((it: any) => {
      it.isChecked = false;
      return status.checklist;
    });
    // for bydefault fetchApi after clearAll checkbox
    fetchApi('all', 'false', ' ');
    setClear('');
  };
  let resourceType: any;
  if (status !== undefined) {
    const resource = status.checklist.slice(0, 2);
    resourceType = resource.map((it: any, idx: number) => (
      <div key={`res-${idx}`} style={{marginBottom: '0.5em'}}>
        <Checkbox
          onClick={filterApi}
          isChecked={checkBoxStatus[it.value]}
          style={{width: '1.2em', height: '1.2em'}}
          label={it.value[0].toUpperCase() + it.value.slice(1)}
          value={it.value}
          name="type"
          id={it.id}
          aria-label="uncontrolled checkbox example"

        />
      </div>
    ));
  }
  let showverifiedtask: any;
  // jsx element for show verifiedtask
  if (status !== undefined) {
    const verifiedtask = status.checklist.slice(2, 3);
    showverifiedtask = verifiedtask.map((it: any, idx: number) => (
      <div key={`task-${idx}`} style={{marginBottom: '0.5em'}}>
        <Checkbox
          onClick={filterApi}
          isChecked={checkBoxStatus[it.value]}
          style={{width: '1.2em', height: '1.2em'}}
          label={it.value[0].toUpperCase() + it.value.slice(1)}
          value={it.value}
          name="verified"
          id={it.id}
          aria-label="uncontrolled checkbox example"

        />
      </div>
    ));
  }
  // jsx element for showing all categories
  let categoryList: any = '';
  if (status !== undefined) {
    const tempstatus = status.checklist.slice(3);
    tempstatus.sort((a: any, b: any) =>
      (a.value > b.value) ? 1 :
        ((b.value > a.value) ? -1 : 0));
    categoryList =
      tempstatus.map((it: any) => (
        <div key={it} style={{marginBottom: '0.5em'}}>
          <Checkbox
            onClick={filterApi}
            isChecked={checkBoxStatus[it.value]}
            style={{width: '1.2em', height: '1.2em'}}
            label={it.value[0].toUpperCase() + it.value.slice(1)}
            value={it.value}
            name="Tags"
            id={it.id}
            aria-label="uncontrolled checkbox example"

          />
        </div>
      ));
  }


  return (
    <div className="filter-size" >
      <h2 style={{marginBottom: '1em'}}>
        {' '}
        <Button component='a' variant='link'
          onClick={clearFilter}> {clear} </Button>
        {'  '}

      </h2>
      <h2 style={{marginBottom: '1em'}}>
        {' '}
        <b>Types</b>{'  '}
      </h2>
      {resourceType}
      <h2 style={{marginBottom: '1em', marginTop: '1em'}}>
        {' '}
        <b>Verified </b>{'  '}
        <Tooltip content={<div>
          Verified Task and Pipelines by Tekton Catalog</div>}>
          <InfoCircleIcon />
        </Tooltip>
      </h2>
      {showverifiedtask}
      <h2 style={{marginBottom: '1em', marginTop: '1em'}}><b>Categories</b></h2>
      {categoryList}
    </div>
  );
};
export default Filter;
