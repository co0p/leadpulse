export namespace employees {
	
	export class Employee {
	    id: number;
	    firstName: string;
	    secondName: string;
	    seniority: string;
	    startDate: string;
	
	    static createFrom(source: any = {}) {
	        return new Employee(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.firstName = source["firstName"];
	        this.secondName = source["secondName"];
	        this.seniority = source["seniority"];
	        this.startDate = source["startDate"];
	    }
	}
	export class Input {
	    firstName: string;
	    secondName: string;
	    seniority: string;
	    startDate: string;
	
	    static createFrom(source: any = {}) {
	        return new Input(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.firstName = source["firstName"];
	        this.secondName = source["secondName"];
	        this.seniority = source["seniority"];
	        this.startDate = source["startDate"];
	    }
	}

}

export namespace overview {
	
	export class Anomaly {
	    employeeId: number;
	    name: string;
	    type: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new Anomaly(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.employeeId = source["employeeId"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.message = source["message"];
	    }
	}
	export class EmployeeHealth {
	    employeeId: number;
	    name: string;
	    morale: number;
	    execution: number;
	    impact: number;
	    multiplier: number;
	    growth: number;
	    culture: number;
	    projectCount: number;
	    evidenceCount: number;
	
	    static createFrom(source: any = {}) {
	        return new EmployeeHealth(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.employeeId = source["employeeId"];
	        this.name = source["name"];
	        this.morale = source["morale"];
	        this.execution = source["execution"];
	        this.impact = source["impact"];
	        this.multiplier = source["multiplier"];
	        this.growth = source["growth"];
	        this.culture = source["culture"];
	        this.projectCount = source["projectCount"];
	        this.evidenceCount = source["evidenceCount"];
	    }
	}
	export class EvidenceGap {
	    employeeId: number;
	    name: string;
	    metric: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new EvidenceGap(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.employeeId = source["employeeId"];
	        this.name = source["name"];
	        this.metric = source["metric"];
	        this.message = source["message"];
	    }
	}
	export class HealthDimensions {
	    execution: number;
	    impact: number;
	    multiplier: number;
	    growth: number;
	    culture: number;
	
	    static createFrom(source: any = {}) {
	        return new HealthDimensions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.execution = source["execution"];
	        this.impact = source["impact"];
	        this.multiplier = source["multiplier"];
	        this.growth = source["growth"];
	        this.culture = source["culture"];
	    }
	}
	export class Overview {
	    headcount: number;
	    rollingTeamMorale: number;
	    highCapacityCount: number;
	    teamMultiplier: number;
	    healthDimensions: HealthDimensions;
	    employeeHealth: EmployeeHealth[];
	    anomalies: Anomaly[];
	    evidenceGaps: EvidenceGap[];
	
	    static createFrom(source: any = {}) {
	        return new Overview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.headcount = source["headcount"];
	        this.rollingTeamMorale = source["rollingTeamMorale"];
	        this.highCapacityCount = source["highCapacityCount"];
	        this.teamMultiplier = source["teamMultiplier"];
	        this.healthDimensions = this.convertValues(source["healthDimensions"], HealthDimensions);
	        this.employeeHealth = this.convertValues(source["employeeHealth"], EmployeeHealth);
	        this.anomalies = this.convertValues(source["anomalies"], Anomaly);
	        this.evidenceGaps = this.convertValues(source["evidenceGaps"], EvidenceGap);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

