import axios from "axios";

const comp = () => {
    const [posts, setPosts] = useState([]);
    const [isLoad, setIsLoad] = useState(false);
    const [isFollow, setIsFollow] = useState(false);

    const follow = () => {
        axios.post('https://dog.ceo/api/breeds/list/all')
            .then(res => res.data)
            .then(data => {
                setIsFollow(!isFollow);
            })
            .catch(error => {
                console.log(error);
            })
    }

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
                    <div className="profile_wrap">
                        <div className="profile_header">
                            <div className="profile_icon_wrap">
                                <img src={profileIconSrc}></img>
                            </div>
                            <div className="profile_name">{username}</div>
                            <div className=`profile_followBtn ${isFollow ? 'active' : ''}` onclick={follow}>
                                {isFollow ? '팔로우' : '언팔로우'}
                            </div>
                            <div className="profile_posts_wrap">
                                <i>포스트 수</i>
                                <div>{postCount}</div>
                            </div>
                            <div className="profile_follower_wrap">
                                <i>팔로워 수</i>
                                <div>{followerCount}</div>
                            </div>
                            <div className="profile_following_wrap">
                                <i>팔로잉 수</i>
                                <div>{followingCount}</div>
                            </div>
                        </div>

                        {
                            posts.map(() => {
                                return Post();
                            })
                        }
                    </div>
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