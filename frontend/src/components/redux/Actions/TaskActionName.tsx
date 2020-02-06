import {FETCH_TASK_NAME} from '../Actions/TaskActionType';
import {API_URL} from '../../../constants';

// eslint-disable-next-line require-jsdoc
export function fetchTaskName(id: number) {
  return function(dispatch: any) {
    fetch(`${API_URL}/resource/${id}`)
        .then((response) => response.json())
        .then((TaskName) => dispatch({
          type: FETCH_TASK_NAME,
          payload: TaskName,
        }));
  };
}

export default fetchTaskName;
