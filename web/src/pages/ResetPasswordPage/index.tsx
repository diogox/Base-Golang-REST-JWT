import React, { Component } from 'react'
import { Link } from 'react-router-dom'
import {login, resetPassword} from '../../utils/AuthService'

export default class ResetPasswordPage extends Component<any> {
    state = {
        newPassword: '',
        errorMsg: '',
        successMsg: '',
    };

    handleSubmit = (ev: React.SyntheticEvent<any>) => {
        ev.preventDefault();

        let token = this.props.match.params.token;

        resetPassword(token, this.state.newPassword)
            .then((res: Response) =>{
                this.setState({
                    successMsg: 'Password successfully changed!',
                });
            })
            .catch((err: Error) =>{
                this.setState({
                    errorMsg: err.message,
                });
            })
    };

    render() {
        const formStyles = "bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4";
        const labelStyles = "block text-grey-darker text-sm font-bold mb-2";
        const inputStyles = "shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline bg-grey-lightest";
        const btnStyles = "bg-blue hover:bg-blue-dark text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline";

        return (
            <div className="container mx-auto m-5">
                <div className="flex justify-center">
                    <form className={formStyles} onSubmit={this.handleSubmit}>
                        <div className="mb-6">
                            <label className={labelStyles}>
                                New Password
                            </label>

                            {/* New Password Input */}
                            <input
                                className={inputStyles}
                                type="password"
                                value={this.state.newPassword}
                                onChange={(ev) => this.setState({ newPassword: ev.target.value })} />

                            {/* Error Message */}
                            {
                                this.state.errorMsg !== '' ?
                                    <p className="text-red text-xs italic">{this.state.errorMsg}</p>
                                    : undefined
                            }

                            {/* Success Message */}
                            {
                                this.state.successMsg !== '' ?
                                    <p className="text-green text-xs italic">{this.state.successMsg}</p>
                                    : undefined
                            }
                        </div>
                        <div className="flex items-center justify-between">
                            <button className={btnStyles} type="submit">
                                Set New Password
                            </button>
                            <Link to="/login" className="inline-block align-baseline font-bold text-sm text-blue hover:text-blue-darker">
                                Remember the password?
                            </Link>
                        </div>
                    </form>
                </div>
            </div>
        )
    }
}