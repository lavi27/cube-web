import axios from "axios";

const comp = () => {
    const [posts, setPosts] = useState([]);
    const [isLoad, setIsLoad] = useState(false);

    useEffect(() => {
        axios.get('https://dog.ceo/api/breeds/list/all')
            .then(res => res.data)
            .then(data => {
                setPosts(data.post);
                setIsLoad(true);
            })
            .catch(error => {
                console.log(error);
            })
    }, [])

    return (
        <div className="home_wrap">
            {
                isLoad ? 
                    posts.map(({content, like, date, userName, userIconSrc}) => {
                        return(
                            <div className="post">
                                <div className="post_header">
                                    <div className="user_wrap">
                                        <div className="user_icon_wrap">
                                            <img src={userIconSrc}></img>
                                        </div>
                                        <span className="user_userName">{userName}</span>
                                    </div>
                                    <span className="date">{date}</span>
                                </div>
                                <div className="post_content">{content}</div>
                                <div className="post_footer">
                                    <div className="like_wrap">
                                        <div className="likeImg_wrap">
                                            <img src="heart"></img>
                                        </div>
                                        <span>{like}</span>
                                    </div>
                                </div>
                            </div>
                        )
                    })
                :
                    <div className="post_skeleton">
                        <div className="post">
                            <div className="post_header">
                                <div className="user_wrap">
                                    <div className="user_icon_wrap">
                                    </div>
                                    <span className="user_userName"></span>
                                </div>
                            </div>
                            <div className="post_content"></div>
                            <div className="post_footer">
                                <div className="like_wrap">
                                </div>
                            </div>
                        </div>
                        <div className="post">
                            <div className="post_header">
                                <div className="user_wrap">
                                    <div className="user_icon_wrap">
                                    </div>
                                    <span className="user_userName"></span>
                                </div>
                            </div>
                            <div className="post_content"></div>
                            <div className="post_footer">
                                <div className="like_wrap">
                                </div>
                            </div>
                        </div>
                    </div>
            }
        </div>
    )
}

export default comp;