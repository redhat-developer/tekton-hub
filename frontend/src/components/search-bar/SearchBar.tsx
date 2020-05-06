/* eslint-disable max-len */
import React, {useState} from 'react';
import {connect} from 'react-redux';
import '@patternfly/react-core/dist/styles/base.css';
import './index.css';
import {
  InputGroup,
  Dropdown,
  DropdownToggle,
  DropdownItem,
  Flex,
  TextInput,
  FlexItem,
} from '@patternfly/react-core';
import {fetchTaskSuccess} from '../redux/Actions/TaskAction';
import {fetchTaskList} from '../redux/Actions/TaskDataListAction';
import store from '../redux/store';
export interface TaskPropData {
  id: number;
  name: string,
  description: string,
  rating: number,
  tags: [],
  lastUpdatedAt: string;
  latestVersion;
  catalog: [],
  type: string,
}

const SearchBar: React.FC = (props: any) => {
  const [sort, setSort] = useState('Name');
  let tempArr: any = [];
  React.useEffect(() => {
    props.fetchTaskSuccess();
    props.fetchTaskList();
    // eslint-disable-next-line
  }, []);
  // Getting all data from store
  if (props.TaskData) {
    tempArr = props.TaskData.map((task: any) => {
      const taskData: TaskPropData = {
        id: task.id,
        catalog: task.catalog,
        name: task.name,
        description: task.description,
        rating: task.rating,
        tags: task.tags,
        type: task.type,
        lastUpdatedAt: task.lastUpdatedAt,
        latestVersion: task.latestVersion,
      };
      return taskData;
    });
  }

  // Dropdown menu
  const [isOpen, set] = useState(false);
  const dropdownItems = [
    <DropdownItem key="name" onClick={sortByName}>Name</DropdownItem>,
    <DropdownItem key="Rating" onClick={sortByRatings}>Ratings</DropdownItem>,
  ];
  const ontoggle = (isOpen: React.SetStateAction<boolean>) => set(isOpen);
  const onSelect = () => set(!isOpen);

  // eslint-disable-next-line require-jsdoc
  function sortByName(event: any) {
    setSort(event.target.text);
    const taskarr = tempArr.sort((first: any, second: any) => {
      if (first.name > second.name) {
        return 1;
      } else {
        return -1;
      }
    });
    store.dispatch({type: 'FETCH_TASK_SUCCESS', payload: taskarr});
  }

  // eslint-disable-next-line require-jsdoc
  function sortByRatings(event: any) {
    setSort(event.target.text);
    const taskarr = tempArr.sort((first: any, second: any) => {
      if (first.rating < second.rating) {
        return 1;
      } else {
        return -1;
      }
    });
    store.dispatch({type: 'FETCH_TASK_SUCCESS', payload: taskarr});
  }

  // AutoComplete text while searching a task
  const [text, setText] = React.useState('');

  const onTextChanged = (e: any) => {
    const value = e;
    let suggestions: any = [];
    const regex = new RegExp(`${value}`, 'i');
    suggestions = props.TaskDataList.sort().filter((v: any) => regex.test(v.name));
    if (value.length === 0) {
      store.dispatch({
        type: 'FETCH_TASK_SUCCESS', payload: props.TaskDataList.sort((first: any, second: any) =>
          first.name > second.name ? 1 : -1),
      });
    } else {
      store.dispatch({
        type: 'FETCH_TASK_SUCCESS', payload: suggestions.sort((first: any, second: any) =>
          first.name > second.name ? 1 : -1),
      });
    }
    setText(value);
  };

  const textValue = text;
  store.dispatch({
    type: 'SEARCH_TEXT',
    payload: textValue,
  });

  return (

    <div className="search">
      <Flex breakpointMods={[{modifier: 'flex-1', breakpoint: 'lg'}]}>

        <Flex breakpointMods={[{modifier: 'column', breakpoint: 'lg'}]}>

          <FlexItem>
            <h2 >
              {' '}
              <b>Sort</b>{'  '}
            </h2>
          </FlexItem>

          <FlexItem>
            <div className="filter">
              <Dropdown
                onSelect={onSelect}
                toggle={<DropdownToggle onToggle={ontoggle}>{sort}</DropdownToggle>}
                isOpen={isOpen}
                dropdownItems={dropdownItems}
              />
            </div>
          </FlexItem>
        </Flex>
        <React.Fragment>


          <InputGroup style={{width: '70%', marginLeft: '11em'}}>
            <div style={{width: '100%', boxShadow: 'rgba'}}>
              <TextInput aria-label="text input example" value={textValue} type="search"
                onChange={onTextChanged} placeholder="Search for task or pipeline"
                style={{padding: '10px 5px'}} />

            </div>
          </InputGroup>


        </React.Fragment>
      </Flex>
    </div>
  );
};

const mapStateToProps = (state: any) => ({
  TaskData: state.TaskData.TaskData,
  TaskDataList: state.TaskDataList.TaskDataList,
  ResourceList: state.ResourceList.ResourceList,

});

export default connect(mapStateToProps, {fetchTaskSuccess, fetchTaskList})(SearchBar);


