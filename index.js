Vue.createApp({
    data() {
        return {
            articles: [],
        }
    },
    mounted() {
        axios.get('https://raw.githubusercontent.com/s-takehana/tech-blog-trend/main/zenn.json'
        ).then(
            response => this.articles = response.data.articles
        ).catch(
            error => console.log(error)
        )
    }
}).mount('#zenn');

Vue.createApp({
    data() {
        return {
            articles: [],
        }
    },
    mounted() {
        axios.get('https://raw.githubusercontent.com/s-takehana/tech-blog-trend/main/qiita.json'
        ).then(
            response => this.articles = response.data
        ).catch(
            error => console.log(error)
        )
    }
}).mount('#qiita');

Vue.createApp({
    data() {
        return {
            articles: [],
        }
    },
    mounted() {
        axios.get('https://dev.to/api/articles', {
            params: {
                top: 1,
            }
        }).then(
            response => this.articles = response.data
        ).catch(
            error => console.log(error)
        )
    }
}).mount('#dev');
