import {combineReducers} from 'redux';
import TaskReducer from './TaskReducer';
import TaskReducerName from './TaskReducerName';
import TaskReducerDescription from './TaskReducerDescription';
import CheckAuthentication from './CheckAuthentication';
import TaskDataListReducer from './TaskDataListReducer';
export default combineReducers({
  TaskData: TaskReducer,
  TaskDataList: TaskDataListReducer,
  TaskName: TaskReducerName,
  TaskDescription: TaskReducerDescription,
  TaskYaml: TaskReducerDescription,
  isAuthenticated: CheckAuthentication,
});


