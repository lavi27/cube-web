const comp = () => {
    const navList = [
        {dir: "/", name="홈", iconSrc: ""},
        {dir: "/", name="프로필", iconSrc: ""},
    ]

    return (
        <div className={style.nav_wrap}>
            {
                navList.map(({dir, name, iconSrc}) => {
                    <Link href={dir}>
                        <div className={style.nav_item}>
                            <img src={iconSrc}></img>
                            <div>{name}</div>
                        </div>
                    </Link>
                })
            }
        </div>
    )
}

export default comp;