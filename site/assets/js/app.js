var app = new Vue({
    el: '#app',
    data: {
        phone: '',
        udata: {},
        searchDone: false
    },
    methods: {
        searchPhone: function () {
            this.$http.get('http://127.0.0.1:6060/phones/' + this.phone).then(
                function (data) {
                    this.udata = data.data;
                    this.searchDone = true;
                }, function (error) {
                    this.udata = {};
                    this.searchDone = true;
                });
        }
    }
});