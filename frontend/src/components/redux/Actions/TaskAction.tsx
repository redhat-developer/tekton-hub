import {FETCH_TASK_SUCCESS} from '../Actions/TaskActionType';
import {API_URL} from '../../../constants';

// eslint-disable-next-line require-jsdoc
export function fetchTaskSuccess() {
  return function(dispatch: any) {
    fetch(`${API_URL}/resources`)
        .then((response) => response.json())
        .then((TaskData) => dispatch({
          type: FETCH_TASK_SUCCESS,
          payload: TaskData,
        }));
  };
}

export default fetchTaskSuccess
;
