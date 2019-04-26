import { combineReducers } from 'redux';
import { userReducer } from "./user_reducer";

export const initialState = {
    //articles: []
    user: {
        username: '',
        email: '',
    }
};

const rootReducer = combineReducers({
    user: userReducer,
});

export default rootReducer;