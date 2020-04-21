/* eslint-disable max-len */
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { SearchIcon, FilterIcon } from '@patternfly/react-icons';
import '@patternfly/react-core/dist/styles/base.css';
import './index.css';
import {
  Button,
  ButtonVariant,
  InputGroup,
  Dropdown,
  DropdownToggle,
  DropdownItem,
  Flex,
  TextInput,
  Card,
} from '@patternfly/react-core';
import { fetchTaskSuccess } from '../redux/Actions/TaskAction';
import { fetchTaskList } from '../redux/Actions/TaskDataListAction';
import store from '../redux/store';
export interface TaskPropData {
  id: number;
  name: string,
  description: string,
  rating: number,
  yaml: string,
  tags: [],
  verified: boolean,
  type: string,
}

const SearchBar: React.FC = (props: any) => {
  const [sort, setSort] = useState('Sort');
  let tempArr: any = [];
  const tempTask: any = [];
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
        name: task.name,
        description: task.description,
        rating: task.rating,
        yaml: task.yaml,
        tags: task.tags,
        verified: task.verified,
        type: task.type,
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
    store.dispatch({ type: 'FETCH_TASK_SUCCESS', payload: taskarr });
  }

  // eslint-disable-next-line require-jsdoc

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
    store.dispatch({ type: 'FETCH_TASK_SUCCESS', payload: taskarr });
  }

  // Searching a task
  let [tasks, setTasks] = useState(''); // Get the user input

  // Search a task
  const searchTask = (text: string) => {
    const task = {
      text,
    };
    tasks = task.text; // user input
    setTasks(tasks);
    if (props.TaskData != null) {
      for (let i = 0; i < tempArr.length; i++) {
        if (tasks.toLocaleLowerCase() === tempArr[i].name.toLocaleLowerCase()) {
          tempTask.push(tempArr[i]);
          break;
        }
      }
    }

    if (tempTask.length > 0) {
      store.dispatch({ type: 'FETCH_TASK_SUCCESS', payload: tempTask });
    }
  };
  const taskNameArr: any = [];
  if (props.TaskDataList) {
    for (let i = 0; i < tempArr.length; i++) {
      taskNameArr.push(tempArr[i].name);
    }
  }

  // AutoComplete text while searching a task
  const [suggestions, setState] = React.useState([]);
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

  // Suggesting a task name
  const suggestionSelected = (value: any) => {
    setText(value);
    setState([]);
    searchTask(value);
  };

  const textValue = text;
  return (

    <div className="search">
      <Flex breakpointMods={[{ modifier: 'flex-1', breakpoint: 'lg' }]}>
        <React.Fragment>

          <InputGroup style={{ width: '70%', marginLeft: '1m' }}>
            <div style={{ width: '100%', boxShadow: 'rgba' }}>
              <TextInput aria-label="search-box" value={textValue} type="search"
                onChange={onTextChanged} placeholder="Search for task or pipeline"
                style={{ outline: 'none', boxSizing: 'border-box', padding: '10px 5px' }} />

              <div style={{ position: 'relative' }}>
                <ul
                  style={{ textAlign: 'left', backgroundColor: 'white', margin: 0, position: 'absolute', width: '100%' }}
                >
                  {suggestions.map((item: any, index: any) =>

                    <li
                      style={{
                        listStyle: 'none', textAlign: 'left', cursor: 'pointer', padding: '10px 7px',
                        border: '0.01em solid rgb(224,224,224)', borderTop: '0',
                      }}
                      onClick={() => suggestionSelected(item)} key={index}>
                      {item}
                    </li>,

                  )}
                </ul>
              </div>

            </div>
          </InputGroup>
          <Card style={{ marginLeft: '-1em' }}>
            <Button variant={ButtonVariant.control} aria-label="search button for search input" >
              <SearchIcon />
            </Button>
          </Card>

          <div className="filter">
            <Button variant={ButtonVariant.control} aria-label="search button for search input" >
              <FilterIcon />
            </Button>
            <Dropdown
              onSelect={onSelect}
              toggle={<DropdownToggle onToggle={ontoggle}>{sort}</DropdownToggle>}
              isOpen={isOpen}
              dropdownItems={dropdownItems}
            />
          </div>
        </React.Fragment>
      </Flex>
    </div>
  );
};

const mapStateToProps = (state: any) => ({
  TaskData: state.TaskData.TaskData,
  TaskDataList: state.TaskDataList.TaskDataList,
});

export default connect(mapStateToProps, { fetchTaskSuccess, fetchTaskList })(SearchBar);


