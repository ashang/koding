class ActivityInnerNavigation extends CommonInnerNavigation

  viewAppended:()->

    # followers...
    filterController = @setListController
      type : "showme"
      itemClass : ListGroupShowMeItem
    , @followerMenuData
    @addSubView filterController.getView()
    filterController.selectItem filterController.getItemsOrdered()[0]

    # everything...
#    filterController = @setListController
#      type : "showme"
#      itemClass : ListGroupShowMeItem
#    , @showMenuData
#    @addSubView filterController.getView()
#    filterController.selectItem filterController.getItemsOrdered()[0]

    @addSubView helpBox = new HelpBox
      subtitle    : "About Your Activity Feed"
      tooltip     :
        title     : "<p class=\"bigtwipsy\">The Activity feed displays posts from the people and topics you follow on Koding. It's also the central place for sharing updates, code, links, discussions and questions with the community. </p>"
        placement : "above"
        offset    : 0
        delayIn   : 300
        html      : yes
        animate   : yes

  # TODO: Question ?? will we show those "JComment" "JOpinion" "JCodeShare"
  followerMenuData :
    title : "SHOW ME"
    items : [      
        { title : "Public", type: "Public" }
        { title : "Followed", type: "Followed" }
        { title : "Status Updates", type : "JStatusUpdate" }
        { title : "Blog Posts", type : "JBlogPost" }
        { title : "Code Snippets", type : "JCodeSnip" }
        { title : "Discussions", type : "JDiscussion" }
        { title : "Tutorials", type : "JTutorial" }
        { title : "Links", type : "JLink", disabledForBeta : yes }

      ]

#  showMenuData :
#    title : "SHOW ME"
#    items : [
#        { title : "Everything" }
#        { title : "Status Updates",   type : "CStatusActivity" }
#        { title : "Blog Posts",       type : "CBlogPostActivity" }
#        { title : "Code Snippets",    type : "CCodeSnipActivity" }
#        { title : "Discussions",      type : "CDiscussionActivity" }
#        { title : "Tutorials",        type : "CTutorialActivity" }
#        { title : "Links",            type : "CLinkActivity", disabledForBeta : yes }
#        # { title : "Code Shares",      type : "codeshare", disabledForBeta : yes }
#        # { title : "Commits",          type : "commit", disabledForBeta : yes }
#        # { title : "Projects",         type : "newproject", disabledForBeta : yes }
#      ]