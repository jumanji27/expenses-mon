export default class Month extends Backbone.View {
  constructor() {
    super();
  }

  returnHTML(month) {
    let WEEKS_IN_MONTH = 5;

    let monthToTemplate = [],
      prevWeek = 0

    for (let key in month) {
      let weekGap = month[key].week - prevWeek;

      if (weekGap > 1) {
        for (let gapKey = 1; gapKey < weekGap; gapKey++) {
          monthToTemplate.push({
            value: 0
          });
        }
      }

      monthToTemplate.push({
        value: month[key].value
      });

      if (month.length === (parseInt(key) + 1) && month[key].week !== WEEKS_IN_MONTH) {
        for (let lastMonthKey = 1; lastMonthKey <= WEEKS_IN_MONTH - month[key].week; lastMonthKey++) {
          monthToTemplate.push({
            value: 0
          });
        }
      }

      prevWeek = month[key].week;
    }

    return tmpl_components_shared_month_main({
      month: monthToTemplate
    });
  }
}