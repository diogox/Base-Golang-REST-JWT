import React, { Component } from 'react'
import { Link } from 'react-router-dom'
import { sendVerificationEmail } from '../../utils/AuthService'

export default class ResendEmailVerificationPage extends Component<any> {
    state = {
        email: '',
        successMsg: '',
        errorMsg: '',
    }

    handleSubmit = (ev: React.SyntheticEvent<any>) => {
        ev.preventDefault();

        // Resend email verification
        sendVerificationEmail(this.state.email)
            .then((r: JSON) =>{
                this.setState({
                    successMsg: "Verification email sent! Check your email to verify your account!",
                })
            })
            .catch((err: Error) =>{
                console.log(err)
                this.setState({
                    errorMsg: err.message,
                })
            })
    }

    render() {
        const formStyles = "bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4";
        const labelStyles = "block text-grey-darker text-sm font-bold mb-2";
        const inputStyles = "shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline";
        const btnStyles = "bg-blue hover:bg-blue-dark text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline";

        return (
            <div className="container mx-auto m-5">
                <div className="flex justify-center">
                    <form className={formStyles} onSubmit={this.handleSubmit}>
                        <div className="mb-4">
                            <label className={labelStyles}>
                                Email
                            </label>

                            {/* Username Input */}
                            <input 
                                className={inputStyles} 
                                value={this.state.email} 
                                onChange={(ev) => this.setState({ email: ev.target.value })} />

                            {
                                this.state.successMsg === '' ? undefined : <p>{this.state.successMsg}</p>
                            }
                            {
                                this.state.errorMsg === '' ? undefined : <p>{this.state.errorMsg}</p>
                            }
                        </div>
                        <div className="flex items-center justify-between">
                            <button className={btnStyles} type="submit">
                                Resend Verification
                            </button>
                            <Link to="signup" className="inline-block align-baseline font-bold text-sm text-blue hover:text-blue-darker">
                                Don't have an account?
                            </Link>
                        </div>
                    </form>
                </div>
            </div>
        )
    }
}