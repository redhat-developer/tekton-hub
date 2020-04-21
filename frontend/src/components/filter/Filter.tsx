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
} from '@patternfly/react-core/dist/js/components';
import {fetchTaskList} from '../redux/Actions/TaskDataListAction';
import {fetchResourceList} from '../redux/Actions/ResourcesList';
import {API_URL} from '../../constants';
import {InfoCircleIcon, DomainIcon, BuildIcon} from '@patternfly/react-icons';
import store from '../redux/store';
import {
  FETCH_TASK_SUCCESS,
  FETCH_TASK_LIST,
}
  from '../redux/Actions/TaskActionType';
import './filter.css';
import {FlexModifiers, Flex, FlexItem} from '@patternfly/react-core';
import {connect} from 'react-redux';
const tempObj: any = {};
const Filter: React.FC = (props: any) => {
  const [categoriesList, setCategoriesList] = useState();
  const [status, setStatus] = useState({checklist: []});
  // const [clear, setClear] = useState(' ');
  const filterItem: any = [{id: '1000', value: 'task', isChecked: false},
    {id: '1001', value: 'pipeline', isChecked: false},
    {id: '1002', value: 'official', isChecked: false},
    {id: '1003', value: 'verified', isChecked: false},
    {id: '1004', value: 'community', isChecked: false}];
  const [checkBoxStatus, setCheckBoxStatus] = useState(
      {},
  );
  //  function for adding categories to filteritem
  const addCategory = (categoryData: any) => {
    categoryData.map((categoryName: string, index: number) =>
      filterItem.push(
          {
            id: `${categoryName['id']}`,
            value: categoryName['name'], isChecked: false,
          },
      ));
    setStatus({checklist: filterItem});
    return categoryData;
  };
  useEffect(() => {
    fetchResourceList();
    fetchTaskList();
    fetch(`https://api-tekton-hub.apps.cluster-blr-8fcf.blr-8fcf.example.opentlc.com/categories`)
        .then((res) => res.json())
        .then((categoryData) =>
          setCategoriesList(addCategory(categoryData.data)));
    if (categoriesList) {
      (Object.keys(categoriesList)).map((category) => {
        return tempObj.category = false;
      });
    }
    setCheckBoxStatus(tempObj);

    // eslint-disable-next-line
  }, []);
  // console.log("categoriesList", categoriesList)
  // console.log("checklist", status.checklist)

  // fetch api function
  const fetchApi = (typeurl: string, verifiedurl: string, categoryurl: any) => {
    fetch(`${API_URL}/resources/${typeurl}/${verifiedurl}?tags=${categoryurl} `)
        .then((resp) => resp.json())
        .then((data) => {
          const taskarr = data.sort((a: any, b: any) =>
          a.name > b.name ? 1 : -1);
          store.dispatch(
              {
                type: FETCH_TASK_SUCCESS,
                payload: taskarr,
              },
          );
          store.dispatch(
              {
                type: FETCH_TASK_LIST,
                payload: taskarr,
              },
          );
        });
  };
  // / function for showing types
  let typeIcon: any;
  const addIcon = (it: any, idx: number) => {
    typeIcon = idx === 0 ? <BuildIcon
      size="sm" color="black"
      style={{marginLeft: '-0.5em'}} /> :
      <DomainIcon size="sm"
        color="black"
        style={{marginLeft: '-0.5em'}} />;
  };


  // custom label for type filter
  const customLabel = (typeName: string) => {
    return <Flex>
      <FlexItem breakpointMods={[{modifier: FlexModifiers['spacer-xs']}]}>
        {typeIcon}
      </FlexItem>
      <FlexItem>
        {typeName}
      </FlexItem>
    </Flex>;
  };


  // formation of filter url  for calling filterAPi to
  //  fetching task and pipelines

  const filterApi = (event: any) => {
    // console.log("filterstatus", event.target.checked)
    // console.log("names", event.target.value)
    // console.log(status.checklist)
    // console.log("categoriesList", categoriesList)

    const tagsList: any = [];
    // const filterdataList: any = [];
    // const filterResourceList: any = [];

    const resourcetypeList: any = [];
    const resourceVerificationList: any = [];
    const target = event.target;
    // for handling isChecked parameter of checkbox
    setCheckBoxStatus({...checkBoxStatus, [target.value]: target.checked});
    status.checklist.forEach((it: any) => {
      if (it.id === event.target.id) {
        return it.isChecked = event.target.checked;
      }
    },
    );

    status.checklist.slice(0, 2).forEach((item: any) => {
      if (item.isChecked === true) {
        resourcetypeList.push(item.value);
      }
    });
    status.checklist.slice(2, 5).forEach((item: any) => {
      if (item.isChecked === true) {
        resourceVerificationList.push(item.value);
      }
    });


    status.checklist.slice(5).forEach((item: any) => {
      if (item.isChecked === true) {
        categoriesList.map((categorytagList: any) => {
          if (categorytagList.name === item.value) {
            categorytagList.tags.forEach((tags: any) =>
              tagsList.push(tags.name));
          }
        });
      }
    });
    // console.log("taglist", tagsList)
    // console.log("props", props.ResourceList)
    // console.log("type", resourcetypeList)
    // console.log("verfied", resourceVerificationList)

    // console.log("store-objects", store.getState())
    //   return tempcategory;
    // });

    // // for displaying clear filter options
    // let flag: any = false;
    // status.checklist.forEach((it: any) => {
    //   if (it.isChecked === true) {
    //     flag = true;
    //   }
    // });
    // if (flag === true) {
    //   setClear('Clear All');
    // } else {
    //   setClear(' ');
    // }
    // fetchApi(typeurl, verifiedurl, categoryurl);
  };

  //   function for clearing all checkbox
  // const clearFilter = () => {
  //   setCheckBoxStatus(
  //     tempObj,
  //   );
  //   status.checklist.forEach((it: any) => {
  //     it.isChecked = false;
  //   });
  //   // for bydefault fetchApi after clearAll checkbox
  //   fetchApi('all', 'false', ' ');
  //   setClear('');
  // };

  let resourceType: any;
  if (status !== undefined && checkBoxStatus !== undefined) {
    const resource = status.checklist.slice(0, 2);
    resourceType = resource.map((it: any, idx: number) => (
      <div key={`res-${idx}`} style={{marginBottom: '0.5em'}}>
        {addIcon(it, idx)}
        <Checkbox
          onClick={filterApi}
          isChecked={checkBoxStatus[it.value]}
          style={{width: '1.2em', height: '1.2em', marginRight: '.3em'}}
          label={customLabel(it.value[0].toUpperCase() + it.value.slice(1))}
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
  if (status !== undefined && checkBoxStatus !== undefined) {
    const verifiedtask = status.checklist.slice(2, 5);
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
  if (status !== undefined && checkBoxStatus !== undefined) {
    const tempstatus = status.checklist.slice(5);
    tempstatus.sort((a: any, b: any) =>
      (a.value > b.value) ? 1 :
        ((b.value > a.value) ? -1 : 0));
    categoryList =
      tempstatus.map((it: any, idx: number) => (
        <div key={`cat-${idx}`} style={{marginBottom: '0.5em'}}>
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
    <div className="filter-size">

      <h2 style={{marginBottom: '1em'}}>
        {' '}
        {/* <Button component='a' variant='link'
          onClick={clearFilter}>
          {clear} </Button> */}
        {'  '}

      </h2>
      <h2 style={{marginBottom: '1em'}}>
        {' '}
        <b>Types</b>{'  '}
      </h2>
      {resourceType}
      <h2 style={{marginBottom: '1em', marginTop: '1em'}}>
        {' '}
        <b>Verification </b>{'  '}
        <Tooltip content={<div>
          Verification Status Task and Pipelines</div>}>
          <InfoCircleIcon />
        </Tooltip>
      </h2>
      {showverifiedtask}
      <h2 style={{marginBottom: '1em', marginTop: '1em'}}><b>Categories</b></h2>
      {categoryList}
    </div>
  );
};

const mapStateToProps = (state: any) => ({
  ResourceList: state.ResourceList.ResourceList,
  TaskDataList: state.TaskDataList.TaskDataList,
  TaskData: state.TaskData.TaskData,
});
export default
connect(mapStateToProps, fetchResourceList, fetchTaskList)(Filter);

