import React, { Component } from 'react'
import { Link } from 'react-router-dom'
import { sendResetPasswordEmail } from '../../utils/AuthService'

export default class SendResetPasswordEmailPage extends Component<any> {
    state = {
        email: '',
        errorMsg: '',
        successMsg: '',
    };

    handleSubmit = (ev: React.SyntheticEvent<any>) => {
        ev.preventDefault();

        sendResetPasswordEmail(this.state.email)
            .then((res: Response) =>{
                this.setState({
                    successMsg: 'Reset password email sent!',
                })
            })
            .catch((err: Error) =>{
                this.setState({
                    errorMsg: err.message,
                })
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
                                Email
                            </label>

                            {/* Email Input */}
                            <input 
                                className={inputStyles} 
                                type="email" 
                                value={this.state.email} 
                                onChange={(ev) => this.setState({ email: ev.target.value })} />

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
                                Reset Password
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