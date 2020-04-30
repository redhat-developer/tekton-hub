import {FETCH_TASK_SUCCESS} from '../Actions/TaskActionType';
import {API_URL} from '../../../constants';

// eslint-disable-next-line require-jsdoc
export function fetchTaskSuccess() {
  return function (dispatch: any) {

    fetch(`${API_URL}/resources`)
      .then((response) => response.json())
      .then((TaskData) => {
        console.log(TaskData)
        dispatch({
          type: FETCH_TASK_SUCCESS,
          payload: TaskData.data.sort((first: any, second: any) =>
            first.name > second.name ? 1 : -1),
        });
      });
  };
}

export default fetchTaskSuccess
  ;
