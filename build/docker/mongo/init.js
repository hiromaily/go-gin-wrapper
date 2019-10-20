//-----------------------------------------------------------------------------
// For initial setting on MongoDB
//-----------------------------------------------------------------------------
//
var conn;
var host = "127.0.0.1"
//init
function init(){
    print("[function] init();");

    //var port = 27017;
    //printjson(port);
    print("port is " + port + ".");

    try{
        conn = new Mongo(host + ":" + port);
    } catch (err){
        print(err.name + ': "' + err.message +  '" occurred when new Mongo().');
        return false;
    }
    return true;
}


// Create Users On Admin Database
function createUserOnAdmin(){
    print("[function] createUserOnAdmin();");

    //1. admin
    //db = connect("127.0.0.1:27017/admin");
    db = conn.getDB("admin");

    // create user who can anything.
    db.createUser({user: "root", pwd: "password", roles: [ "root" ] });

    // create user who has admin authority on all databases.
    db.createUser({user:"admin", pwd:"admin", roles: [{role: "userAdminAnyDatabase", db: "admin"}] });
}

// Create Users On hiromaily Database
function createUserOnHiromaily(){
    print("[function] createUserOnHiromaily();");

    //2. hiromaily
    //db = connect("127.0.0.1:27017/hiromaily");
    db = conn.getDB("hiromaily");

    // create user who has authority on hiromaily database.
    db.createUser({user:"hiromaily", pwd:"12345678", roles: [{role: "userAdmin", db: "hiromaily"}] });

    // auth
    db.auth("hiromaily","12345678");
}

// Insert records in News Collection
function insertNews(){
    print("[function] insertNews();");

    //2. hiromaily
    db = conn.getDB("hiromaily");

    //news collection
    db.createCollection("news");

    //remove
    db.news.remove({});

    var date = new Date();

    //insert
    db.news.insert({
        "news_id": 1,
        "name": "TechCrunch",
        "url": "http://feeds.feedburner.com/TechCrunch/",
        "categories": [
            {
                "name": "startups",
                "url": "http://feeds.feedburner.com/TechCrunch/startups"
            },{
                "name": "Europe",
                "url": "http://feeds.feedburner.com/Techcrunch/europe"
            },{
                "name": "GreenTech",
                "url": "http://feeds.feedburner.com/TechCrunch/greentech"
            }
        ],
        "createdAt": date
    });

    db.news.insert({
        "news_id": 2,
        "name": "MacRumors",
        "url": "http://feeds.macrumors.com/MacRumors-All",
        "categories": [],
        "createdAt": date
    });

    db.news.insert({
        "news_id": 3,
        "name": "CNET",
        "url": "http://www.cnet.com/rss/all/",
        "categories": [],
        "createdAt": date
    });

    db.news.insert({
        "news_id": 4,
        "name": "The Next Web",
        "url": "http://feeds2.feedburner.com/thenextweb",
        "categories": [],
        "createdAt": date
    });
}

//-----------------------------------------------------------------------------
// Main()
//-----------------------------------------------------------------------------
var bRet = init();
if (bRet) {
    createUserOnAdmin();

    createUserOnHiromaily();

    insertNews();
}

