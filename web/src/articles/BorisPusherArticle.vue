<template>
  <div class="article-body">
    <!-- Table of Contents -->
    <nav class="toc">
      <h3 class="toc-title">Table of Contents</h3>
      <ol class="toc-list">
        <li><a href="#introduction" class="toc-link">Introduction</a></li>
        <li><a href="#velocity-verlet" class="toc-link">The Closer Look at Velocity Verlet</a></li>
        <li><a href="#boris-pusher-algorithm" class="toc-link">The Boris Pusher Algorithm</a></li>
        <li><a href="#implementation" class="toc-link">Implementation</a></li>
        <li><a href="#results" class="toc-link">Results</a></li>
        <li><a href="#outcomes-and-learnings" class="toc-link">Outcomes & Learnings</a></li>
      </ol>
    </nav>

    <!-- Introduction -->
    <h2 id="introduction" class="section-heading">Introduction</h2>
    <p>
      I recently released the initial version of the Simulation Catalogue. Along with it came the
      charged particle in a magnetic field simulation, which calculates the trajectory of a charged
      particle moving in a constant magnetic, electric, and gravitational field (<a
        href="/blog/initial-release"
        class="inline-link"
        >more on that here</a
      >).
    </p>

    <p>
      Despite positive initial results, a closer inspection of the results revealed some issues.
    </p>

    <figure class="result-figure">
      <img
        src="/blog/images/002/verlet_x3_trajectory.png"
        alt="X3 Evolution"
        class="result-image"
      />
      <figcaption class="result-caption">
        Figure 1 <br />
        x3 evolution showing instability in x3 plane.
      </figcaption>
    </figure>

    <br />
    <p>
      At first glance this looks exactly as expected. The magnetic field is causing the charged
      particle to oscillate in the <span v-html="renderLatexInline(`x_3`)"></span> plane. However, a
      closer look shows that the amplitude of the oscillation is increasing over time, ultimately
      resulting in an unstable trajectory.
    </p>

    <h2 id="velocity-verlet" class="section-heading">A Closer Look at Velocity Verlet</h2>
    <p>
      Unstable trajectories usually point to an issue with the underlying algorithm. In this case, I
      was using Velocity Verlet, a common "leapfrog" style algorithm used to model a whole host of
      physical systems with great success.
    </p>

    <p>
      Velocity Verlet updates the velocity of the particle in half-steps, and is implemented in the
      following manner:
    </p>

    <ol class="centered-list">
      <li>
        Calculate
        <span
          v-html="renderLatexInline(`v_{t + \\frac{\\Delta t}{2}} = v_t + \\frac{a\\Delta t}{2}`)"
        ></span>
      </li>
      <li>
        Calculate
        <span
          v-html="
            renderLatexInline(`x_{t + \\Delta t} = x_t + v_{t + \\frac{\\Delta t}{2}}\\Delta t`)
          "
        ></span>
      </li>
      <li>
        Calculate
        <span
          v-html="
            renderLatexInline(
              `v_{t + \\Delta t} = v_{t + \\frac{\\Delta t}{2}} + \\frac{a\\Delta t}{2}`,
            )
          "
        ></span>
      </li>
    </ol>

    <p>where the acceleration is re-computed between each velocity update.</p>

    <p>
      Usually, the result is a sort of self-stabilizing trajectory. Instabilities do occur, but they
      self-regulate over the course of the simulation. The particle may veer off its trajectory
      slightly in one iteration of the loop, but it will veer off in the opposite direction in a
      subsequent iteration. Errors average themselves out.
    </p>

    <p>
      However, this clearly isn't happening here. So what's going on? The answer is that Velocity
      Verlet does not work with velocity dependent forces. Step 3 involves computing the force at
      <span v-html="renderLatexInline(`t + \\Delta t`)"></span>. However, this requires knowing
      <span v-html="renderLatexInline(`v_{t + \\Delta t}`)"></span>. If the force depends on the
      velocity, a circular definition is introduced. The only way to compute the force is to use the
      velocity at <span v-html="renderLatexInline(`v_{t + \\frac{\\Delta t}{2}}`)"></span>, which is
      incorrect. This is the source of the instability.
    </p>

    <p>
      It's worth mentioning here that the initial Velocity Verlet implementation did produce
      sensible results, especially in cases where the time step
      <span v-html="renderLatexInline(`\\Delta t`)"></span>
      is small and the initial conditions are relatively stable.
    </p>

    <p>
      Solving the instability required a more sophisticated algorithm that handled the velocity
      dependence explicitly. This is where the Boris Pusher comes into play.
    </p>

    <!-- The Boris Pusher Algorithm -->
    <h2 id="boris-pusher-algorithm" class="section-heading">The Boris Pusher Algorithm</h2>
    <p>
      Handling the velocity dependence requires a more subtle approach. The Boris Pusher is the gold
      standard when it comes to plasma physics simulations. It separates forces into two separate
      batches: those that cause linear acceleration (electric fields, gravity etc) and those that
      cause gyration of the velocity vector (magnetic) field.
    </p>

    <ol class="centered-list">
      <li>
        Use linear forces to compute the first half-kicked velocity vector
        <span v-html="renderLatexInline(`v_{-}`)"></span>
      </li>
      <li>Rotate the half-kicked velocity vector using the magnetic field</li>
      <li>
        Use linear forces to compute the second half-kicked velocity vector
        <span v-html="renderLatexInline(`v_{+}`)"></span>
      </li>
      <li>Calculate final position</li>
    </ol>

    <p>Mathematically, this can be broken down into the following</p>

    <div v-html="renderLatexDisplay(borisAlgorithm)"></div>

    <p>
      There's an important caveat here: the Boris Pusher implemented this way only works because the
      fields (electric, magnetic and gravitational) are constant i.e they do not depend on the
      position of the particle. This allows us to compute the final velocity and position updates at
      the very end.
    </p>

    <p>
      The Boris Pusher can be used for non-constant fields as well. One easy approach is to simply
      evaluate the fields before the Boris Push is started using the current position vector. This
      is often good enough for fields that vary slowly relative to the gyration radius. Another
      slightly more accurate approach is to estimate the half step position
      <span v-html="renderLatexInline(`x(t + \\frac{\\Delta t}{2})`)"></span> and use that to
      evaluate the fields.
    </p>
    <!-- Implementation -->
    <h2 id="implementation" class="section-heading">Implementation</h2>
    <p>
      With the algorithm in place, it's time to write the Fortran implementation. Since this was the
      second integrator that I wrote, I decided to implement a generic integrator interface. This
      allows me to define a whole series of integrators that I can use interchangeably throughout
      the codebase.
    </p>

    <div class="code-block" v-html="highlightCode(interfaceDefinition, 'fortran')"></div>

    <p>
      Lets quickly break down whats going on here. First I define a
      <span v-html="renderMarkdownInline('`integrator_t`')"></span> abstract type, which has a
      <span v-html="renderMarkdownInline('`integrate_step`')"></span> procedure attached.
      <span v-html="renderMarkdownInline('`integrate_step`')"></span> must implement the
      <span v-html="renderMarkdownInline('`proto_step`')"></span> abstract interface. Think of it
      this way: <span v-html="renderMarkdownInline('`integrator_t`')"></span> is the integrator,
      <span v-html="renderMarkdownInline('`integrate_step`')"></span> defines how the integrator
      performs numerical integration.
    </p>

    <p>
      Note that <span v-html="renderMarkdownInline('`integrate_step`')"></span> relies on a few
      custom types that are defined in other modules, mainly
      <span v-html="renderMarkdownInline('`PointParticle`')"></span> and
      <span v-html="renderMarkdownInline('`force_t`')"></span>. I won't cover these in detail here,
      but <span v-html="renderMarkdownInline('`PointParticle`')"></span> is a type that holds the
      basic properties (mass, charge etc) and mechanical state (velocity, position etc) for a point
      particle. <span v-html="renderMarkdownInline('`force_t`')"></span> is another abstract type
      that has a <span v-html="renderMarkdownInline('`get_force`')"></span> function which generates
      the force vector acting on a given
      <span v-html="renderMarkdownInline('`PointParticle`')"></span> instance.
    </p>

    <p>With the interface in place, the Boris implementation can be created.</p>

    <div class="code-block" v-html="highlightCode(borisImplementation, 'fortran')"></div>

    <p>
      <span v-html="renderMarkdownInline('`boris_t`')"></span> is defined as a type that extends our
      <span v-html="renderMarkdownInline('`integrator_t`')"></span> interface. The magnetic field is
      also set as an attribute of the <span v-html="renderMarkdownInline('`boris_t`')"></span> type.
      Next we define a subroutine called
      <span v-html="renderMarkdownInline('`boris_step`')"></span> which implements the algorithm
      outlined above. Note that the signature of the subroutine must match precisely that of the
      <span v-html="renderMarkdownInline('`proto_step`')"></span> abstract interface, otherwise the
      code won't compile. Finally, the
      <span v-html="renderMarkdownInline('`boris_step`')"></span> subroutine is attached to the
      <span v-html="renderMarkdownInline('`boris_t`')"></span> type as its implementation of
      <span v-html="renderMarkdownInline('`proto_step`')"></span>.
    </p>

    <p>
      All thats left is to replace the old verlet velocity code with the new boris step integrator
    </p>

    <div class="code-block" v-html="highlightCode(borisUsage, 'fortran')"></div>

    <p>
      Hopefully now it becomes clear why I went to the trouble of defining the abstract type
      <span v-html="renderMarkdownInline('`integrator_t`')"></span>. I can now swap out the
      <span v-html="renderMarkdownInline('`boris_t`')"></span> integrator for any other type that
      implements <span v-html="renderMarkdownInline('`integrator_t`')"></span> without changing the
      rest of the code.
    </p>
    <!-- Results -->
    <h2 id="results" class="section-heading">Results</h2>
    <p>
      Now comes the relevant part. It's all for nothing if it doesn't stabilize the trajectory. The
      following simulation configuration was used to test the new Boris implementation.
    </p>
    <div class="code-block" v-html="highlightCode(simConfig, 'toml')"></div>
    <p>
      The following figure shows the individual trajectory components for both the verlet and the
      boris algorithm.
    </p>

    <figure class="result-figure">
      <img
        src="/blog/images/002/verlet_vs_boris_comparison.png"
        alt="X3 Evolution"
        class="result-image"
      />
      <figcaption class="result-caption">
        Figure 2 <br />
        Trajectory for Verlet and Boris Algorithm.
      </figcaption>
    </figure>

    <br />
    <p>
      The evolution in <span v-html="renderLatexInline(`x_1`)"></span> and
      <span v-html="renderLatexInline(`x_2`)"></span> was largely unchanged. Closer inspection of
      <span v-html="renderLatexInline(`x_2`)"></span> does reveal a slightly more stable trajectory.
      The Verlet algorithm doesn't break down quite as much in
      <span v-html="renderLatexInline(`x_2`)"></span> as it does in
      <span v-html="renderLatexInline(`x_3`)"></span>, but the Boris implementation does produce a
      smoother path.
    </p>

    <p>
      The difference in the <span v-html="renderLatexInline(`x_3`)"></span> component however is
      drastic. Much like before, the Verlet algorithm almost instantly produces an unstable
      trajectory. The Boris algorithm on the other hand, is completely stable.
    </p>

    <h2 id="outcomes-and-learnings" class="section-heading">Outcomes and Learnings</h2>

    <p>
      It will come as little surprise to anyone with expertise in the field that the Boris Pusher
      Algorithm outperforms Velocity Verlet on every metric. It's definitely a useful tool, and will
      likely become my de facto algorithm for future simulations that involve magnetic fields.
    </p>
    <p>
      The abstract type approach used in the Fortran code for the integrators and forces proved to
      be useful as well. I was able to run the simulation with the Boris and Verlet algorithms
      almost interchangeably to produce the plots for analysis, with very little code changes
      between iterations.
    </p>
    <p>
      There are definitely improvements to be made, especially to the
      <span v-html="renderMarkdownInline('`force_t`')"></span>
      interface, which really needs to be revised to encapsulate a field rather than a single force.
      Updating the Boris implementation to handle time and position dependent fields is also a
      natural evolution of the current code, and will be required for more complicated systems.
    </p>
    <p>
      But overall, I achieved what I set out to do. The Boris algorithm stabilized the trajectory,
      which not only produced better plots, but also drastically saves me on compute. I was able to
      reduce both the time step and the number of iterations by an order of magnitude in the live
      simulation catalogue. The simulation runs 10x faster, and produces better results. A definite
      success.
    </p>
    <p>
      Some of the code provided in this article was a little rushed. If you're interested in seeing
      the implementation thats actually running on this platform, feel free to visit the
      <a
        href="https://github.com/PSauerborn/simulations/blob/master/app/helical_motion_simulation.f90"
        target="_blank"
        class="inline-link"
        >GitHub repo</a
      >.
    </p>
  </div>
</template>

<script setup>
import {
  renderLatexInline,
  renderLatexDisplay,
  highlightCode,
  renderMarkdownInline,
} from './utils.js'

const borisAlgorithm = `\\begin{aligned}
  \\mathbf{a} &= \\frac{\\mathbf{F}_{linear}}{m} \\\\
  \\mathbf{v}_{-} &= \\mathbf{v}(t) + \\frac{\\mathbf{a} \\Delta t}{2} \\\\
  \\mathbf{t} &= \\frac{q \\mathbf{B} \\Delta t}{2m} \\\\
  \\mathbf{s} &= \\frac{2 \\mathbf{t}}{1 + |\\mathbf{t}|^2} \\\\
  \\mathbf{v}^{\\prime} &= \\mathbf{v}_{-} + (\\mathbf{v}_{-} \\times \\mathbf{t}) \\\\
  \\mathbf{v}_{+} &= \\mathbf{v}_{-} + (\\mathbf{v}^{\\prime} \\times \\mathbf{s}) \\\\
  \\mathbf{v}(t + \\Delta t) &= \\mathbf{v}_{+} + \\frac{\\mathbf{a} \\Delta t}{2} \\\\
  \\mathbf{x}(t + \\Delta t) &= \\mathbf{x}(t) + \\mathbf{v}(t + \\Delta t) \\Delta t
\\end{aligned}`

const interfaceDefinition = `module integrators
  implicit none
  private

  type, abstract, public :: integrator_t
    real :: delta_t

  contains
    procedure(step_proto), deferred, pass(this) :: integrate_step

  end type integrator_t

  abstract interface
      subroutine step_proto(this, particle, force)
         import integrator_t
         import PointParticle
         import force_t

         class(integrator_t), intent(in) :: this  !< The integrator instance
         type(PointParticle), intent(inout) :: particle  !< Particle to update
         class(force_t) :: force  !< Force calculator
      end subroutine step_proto
  end interface
end module integrators`

const borisImplementation = `module integrators
  type, extends(integrator_t), public :: boris_t
    real, dimension(3) :: magnetic_field

  contains
    procedure, pass(this) :: integrate_step => boris_step

  end type boris_t

contains
  subroutine boris_step(this, particle, force)
    class(boris_t), intent(in) :: this
    type(PointParticle), intent(inout) :: particle
    class(force_t), intent(in) :: force

    real :: a_linear(3)
    real :: v_minus(3), v_prime(3), v_plus(3), t_vector(3), s_vector(3)

    ! calculate force due to linear components
    a_linear = force%get_force(particle)/particle%mass
    ! do first velocity update
    v_minus = particle%state%velocity + a_linear*(this%delta_t/2)

    ! define tangent and secant vectors
    t_vector = (particle%charge/particle%mass)*this%magnetic_field*(this%delta_t/2)
    s_vector = (2*t_vector)/(1 + dot_product(t_vector, t_vector))

    v_prime = v_minus + cross_product(v_minus, t_vector)
    v_plus = v_minus + cross_product(v_prime, s_vector)

    particle%state%velocity = v_plus + a_linear*(this%delta_t/2)
    particle%state%position = particle%state%position + particle%state%velocity*this%delta_t

  end subroutine boris_step

  function new_boris_integrator(delta_t, magnetic_field) result(integrator)
    real, intent(in) :: delta_t
    real, intent(in) :: magnetic_field(3)
    type(boris_t) :: integrator

    integrator%delta_t = delta_t
    integrator%magnetic_field = magnetic_field

   end function new_boris_integrator

end module integrators`

const borisUsage = `program boris_example
  use integrators
  implicit none

  type(boris_t) :: boris_integrator
  type(lorentz_force_t) :: force
  type(PointParticle) :: particle

  real, dimension(3) :: magnetic_field, electric_field
  real, dimension(3) :: initial_position, initial_velocity
  real :: charge, mass
  real, dimension(:, :), allocatable :: trajectory
  real :: delta_t
  integer :: num_steps, i

  delta_t = 0.01
  num_steps = 5000

  charge = 1.0
  mass = 1.0

  magnetic_field = [1.0, 0.0, 0.0]
  electric_field = [0.0, 2.0, 1.0]

  initial_position = [0.0, 0.0, 0.0]
  initial_velocity = [1.0, 5.0, 5.0]

  ! define force acting on particle
  force = new_lorentz_force_t(electric_field, magnetic_field)

  allocate(trajectory(num_steps, 3))
  ! create new boris integrator
  boris_integrator = new_boris_integrator(delta_t, magnetic_field)

  ! create new point particle instance
  particle = new_point_particle(mass, charge, initial_position, initial_velocity)

  do i = 1, num_steps
    ! call the integrate_step() subroutine with each iteration
    call boris_integrator%integrate_step(particle, force)
    ! update the trajectory
    trajectory(i, :) = particle%state%position
  end do

end program`

const simConfig = `[parameters]
initial_velocity = [1.0, 5.0, 1.0]
initial_position = [0.0, 0.0, 0.0]
magnetic_field = [0.25, 0.0, 0.0]
electric_field = [0.0, 0.0, 0.0]
mass = 0.25
charge = 0.5
delta_t = 0.01
num_steps = 50000

[config]
output_dir = "output"
`
</script>

<style lang="scss" scoped>
.toc {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  padding: 20px 24px;
  margin-bottom: 32px;

  .toc-title {
    font-size: 0.85rem;
    font-weight: 600;
    color: #636366;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    margin: 0 0 12px 0;
  }

  .toc-list {
    margin: 0;
    padding-left: 20px;
    color: #d1d1d6;

    li {
      margin-bottom: 6px;
      line-height: 1.5;
    }

    ul {
      padding-left: 20px;
      margin: 6px 0 0 0;
    }
  }

  .toc-link {
    color: #30d158;
    text-decoration: none;
    font-size: 0.95rem;

    &:hover {
      text-decoration: underline;
    }
  }
}

.section-heading {
  font-size: 1.5rem;
  font-weight: 700;
  color: #f5f5f7;
  margin: 40px 0 16px 0;
  letter-spacing: -0.01em;

  &:first-of-type {
    margin-top: 0;
  }
}

.subsection-heading {
  font-size: 1.2rem;
  font-weight: 600;
  color: #e5e5ea;
  margin: 28px 0 12px 0;
}

.equation {
  margin: 24px 0;
  padding: 20px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  text-align: center;
  overflow-x: auto;
}

.code-block {
  margin: 24px 0;

  :deep(pre) {
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    padding: 20px;
    overflow-x: auto;
    margin: 0;

    code {
      font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
      font-size: 0.875rem;
      line-height: 1.6;
      color: #d1d1d6;
    }
  }
}

.results-gallery {
  display: flex;
  flex-direction: column;
  gap: 24px;
  margin: 24px 0 32px 0;
}

.result-figure {
  margin: 0;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;

  &:hover {
    background: rgba(255, 255, 255, 0.06);
    border-color: rgba(255, 255, 255, 0.15);
    transform: translateY(-2px);
  }
}

.result-image {
  width: 100%;
  height: auto;
  display: block;
}

.result-caption {
  padding: 16px;
  font-size: 0.9rem;
  line-height: 1.5;
  color: #a1a1a6;
  text-align: center;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.centered-list {
  display: inline-block;
  text-align: left;
  margin: 24px 0;
  padding: 20px 40px;

  li {
    margin-bottom: 8px;
    line-height: 1.8;
  }
}

.article-body {
  text-align: center;

  p,
  h2,
  h3,
  nav,
  figure,
  .code-block {
    text-align: left;
  }
}

// GitHub-style inline code
:deep(code) {
  background-color: rgba(110, 118, 129, 0.4);
  padding: 0.2em 0.4em;
  border-radius: 6px;
  font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
  font-size: 0.875em;
  color: #e6edf3;
}

// Don't double-style code inside code blocks
.code-block :deep(code) {
  background-color: transparent;
  padding: 0;
  border-radius: 0;
}
</style>
